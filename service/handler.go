// Copyright (c) nano Author and TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dansen/pud/defaultlog/log"

	"github.com/nats-io/nuid"

	"github.com/dansen/pud/acceptor"
	"github.com/dansen/pud/pipeline"

	"github.com/dansen/pud/agent"
	"github.com/dansen/pud/cluster"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/conn/codec"
	"github.com/dansen/pud/conn/message"
	"github.com/dansen/pud/conn/packet"
	"github.com/dansen/pud/constants"
	pcontext "github.com/dansen/pud/context"
	"github.com/dansen/pud/docgenerator"
	e "github.com/dansen/pud/errors"
	"github.com/dansen/pud/metrics"
	"github.com/dansen/pud/route"
	"github.com/dansen/pud/serialize"
	"github.com/dansen/pud/session"
	"github.com/dansen/pud/timer"
	"github.com/dansen/pud/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	handlerType = "handler"
)

type (
	// HandlerService service
	HandlerService struct {
		baseService
		chLocalProcess   chan unhandledMessage // channel of messages that will be processed locally
		chRemoteProcess  chan unhandledMessage // channel of messages that will be processed remotely
		decoder          codec.PacketDecoder   // binary decoder
		remoteService    *RemoteService
		serializer       serialize.Serializer          // message serializer
		server           *cluster.Server               // server obj
		services         map[string]*component.Service // all registered service
		metricsReporters []metrics.Reporter
		agentFactory     agent.AgentFactory
		handlerPool      *HandlerPool
		handlers         map[string]*component.Handler // all handler method
	}

	unhandledMessage struct {
		ctx   context.Context
		agent agent.Agent
		route *route.Route
		msg   *message.Message
	}
)

// NewHandlerService creates and returns a new handler service
func NewHandlerService(
	packetDecoder codec.PacketDecoder,
	serializer serialize.Serializer,
	localProcessBufferSize int,
	remoteProcessBufferSize int,
	server *cluster.Server,
	remoteService *RemoteService,
	agentFactory agent.AgentFactory,
	metricsReporters []metrics.Reporter,
	handlerHooks *pipeline.HandlerHooks,
	handlerPool *HandlerPool,
) *HandlerService {
	h := &HandlerService{
		services:         make(map[string]*component.Service),
		chLocalProcess:   make(chan unhandledMessage, localProcessBufferSize),
		chRemoteProcess:  make(chan unhandledMessage, remoteProcessBufferSize),
		decoder:          packetDecoder,
		serializer:       serializer,
		server:           server,
		remoteService:    remoteService,
		agentFactory:     agentFactory,
		metricsReporters: metricsReporters,
		handlerPool:      handlerPool,
		handlers:         make(map[string]*component.Handler),
	}

	h.handlerHooks = handlerHooks

	return h
}

// Dispatch message to corresponding logic handler
func (h *HandlerService) Dispatch(thread int) {
	// TODO: This timer is being stopped multiple times, it probably doesn't need to be stopped here
	defer timer.GlobalTicker.Stop()

	for {
		// Calls to remote servers block calls to local server
		select {
		case lm := <-h.chLocalProcess:
			metrics.ReportMessageProcessDelayFromCtx(lm.ctx, h.metricsReporters, "local")
			h.localProcess(lm.ctx, lm.agent, lm.route, lm.msg)

		case rm := <-h.chRemoteProcess:
			metrics.ReportMessageProcessDelayFromCtx(rm.ctx, h.metricsReporters, "remote")
			h.remoteService.remoteProcess(rm.ctx, nil, rm.agent, rm.route, rm.msg)

		case <-timer.GlobalTicker.C: // execute cron task
			timer.Cron()

		case t := <-timer.Manager.ChCreatedTimer: // new Timers
			timer.AddTimer(t)

		case id := <-timer.Manager.ChClosingTimer: // closing Timers
			timer.RemoveTimer(id)
		}
	}
}

// Register registers components
func (h *HandlerService) Register(comp component.Component, opts []component.Option) error {
	s := component.NewService(comp, opts)

	if _, ok := h.services[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all handlers
	h.services[s.Name] = s
	for name, handler := range s.Handlers {
		h.handlerPool.Register(s.Name, name, handler)
	}
	return nil
}

// Handle handles messages from a conn
func (h *HandlerService) Handle(conn acceptor.PlayerConn) {
	// create a client agent and startup write goroutine
	a := h.agentFactory.CreateAgent(conn)

	// startup agent goroutine
	go a.Handle()

	log.Infof("New session established: %s", a.String())

	// guarantee agent related resource is destroyed
	defer func() {
		log.Infof("Session read goroutine exit, SessionID=%v, UID=%v", a.GetSession().ID(), a.GetSession().UID())
		a.GetSession().Close()
	}()

	for {
		msg, err := conn.GetNextMessage()

		if err != nil {
			if err != constants.ErrConnectionClosed {
				// log.Errorf("Error reading next available message: %s", err.Error())
			}

			return
		}

		packets, err := h.decoder.Decode(msg)
		if err != nil {
			log.Errorf("Failed to decode message: %s", err.Error())
			return
		}

		if len(packets) < 1 {
			log.Warnf("Read no packets, data: %v", msg)
			continue
		}

		// process all packet
		for i := range packets {
			if err := h.processPacket(a, packets[i]); err != nil {
				log.Errorf("Failed to process packet: %s", err.Error())
				return
			}
		}
	}
}

func (h *HandlerService) processPacket(a agent.Agent, p *packet.Packet) error {
	switch p.Type {
	case packet.Handshake:
		log.Debug("Received handshake packet")
		if err := a.SendHandshakeResponse(); err != nil {
			log.Errorf("Error sending handshake response: %s", err.Error())
			return err
		}
		log.Debugf("Session handshake Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

		// Parse the json sent with the handshake by the client
		handshakeData := &session.HandshakeData{}
		err := json.Unmarshal(p.Data, handshakeData)
		if err != nil {
			a.SetStatus(constants.StatusClosed)
			return fmt.Errorf("Invalid handshake data. Id=%d", a.GetSession().ID())
		}

		a.GetSession().SetHandshakeData(handshakeData)
		a.SetStatus(constants.StatusHandshake)
		err = a.GetSession().Set(constants.IPVersionKey, a.IPVersion())
		if err != nil {
			log.Warnf("failed to save ip version on session: %q\n", err)
		}

		log.Debug("Successfully saved handshake data")

	case packet.HandshakeAck:
		a.SetStatus(constants.StatusWorking)
		log.Debugf("Receive handshake ACK Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

	case packet.Data:
		if a.GetStatus() < constants.StatusWorking {
			return fmt.Errorf("receive data on socket which is not yet ACK, session will be closed immediately, remote=%s",
				a.RemoteAddr().String())
		}

		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		h.processMessage(a, msg)

	case packet.Heartbeat:
		// expected
	}

	a.SetLastAt()
	return nil
}

func (h *HandlerService) processMessage(a agent.Agent, msg *message.Message) {
	requestID := nuid.New()
	// 处理消息时才创建一个context上下文
	ctx := pcontext.AddToPropagateCtx(context.Background(), constants.StartTimeKey, time.Now().UnixNano())
	ctx = pcontext.AddToPropagateCtx(ctx, constants.RouteKey, msg.Route)
	ctx = pcontext.AddToPropagateCtx(ctx, constants.RequestIDKey, requestID)
	tags := opentracing.Tags{
		"local.id":   h.server.ID,
		"span.kind":  "server",
		"msg.type":   strings.ToLower(msg.Type.String()),
		"user.id":    a.GetSession().UID(),
		"request.id": requestID,
	}
	ctx = tracing.StartSpan(ctx, msg.Route, tags)
	ctx = context.WithValue(ctx, constants.SessionCtxKey, a.GetSession())

	r, err := route.Decode(msg.Route)
	if err != nil {
		log.Errorf("Failed to decode route: %s", err.Error())
		a.AnswerWithError(ctx, msg.ID, e.NewError(err, e.ErrBadRequestCode))
		return
	}

	if r.SvType == "" {
		r.SvType = h.server.Type
	}

	message := unhandledMessage{
		ctx:   ctx,
		agent: a,
		route: r,
		msg:   msg,
	}

	// 通过第一个serverType标识识别是否是本地服务器
	if r.SvType == h.server.Type {
		h.chLocalProcess <- message
	} else {
		// 远程服务器，请求服务发现的逻辑
		if h.remoteService != nil {
			h.chRemoteProcess <- message
		} else {
			log.Warnf("request made to another server type but no remoteService running")
		}
	}
}

func (h *HandlerService) localProcess(ctx context.Context, a agent.Agent, route *route.Route, msg *message.Message) {
	var mid uint
	switch msg.Type {
	case message.Request:
		mid = msg.ID
	case message.Notify:
		mid = 0
	}

	ret, err := h.handlerPool.ProcessHandlerMessage(ctx, route, h.serializer, h.handlerHooks, a.GetSession(), msg.Data, msg.Type, false)
	if msg.Type != message.Notify {
		if err != nil {
			log.Errorf("Failed to process handler message: %s", err.Error())
			a.AnswerWithError(ctx, mid, err)
		} else {
			err := a.GetSession().ResponseMID(ctx, mid, ret)
			if err != nil {
				tracing.FinishSpan(ctx, err)
				metrics.ReportTimingFromCtx(ctx, h.metricsReporters, handlerType, err)
			}
		}
	} else {
		metrics.ReportTimingFromCtx(ctx, h.metricsReporters, handlerType, nil)
		tracing.FinishSpan(ctx, err)
	}
}

// DumpServices outputs all registered services
func (h *HandlerService) DumpServices() {
	handlers := h.handlerPool.GetHandlers()
	for name := range handlers {
		log.Infof("registered handler %s, isRawArg: %v", name, handlers[name].IsRawArg)
	}
}

// Docs returns documentation for handlers
func (h *HandlerService) Docs(getPtrNames bool) (map[string]interface{}, error) {
	if h == nil {
		return map[string]interface{}{}, nil
	}
	return docgenerator.HandlersDocs(h.server.Type, h.services, getPtrNames)
}
