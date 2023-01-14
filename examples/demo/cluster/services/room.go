package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dansen/pud/logger/log"

	"github.com/dansen/pud"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/examples/demo/protos"
	"github.com/dansen/pud/timer"
	"github.com/google/uuid"
)

type (
	// Room represents a component that contains a bundle of room related handler
	// like Join/Message
	Room struct {
		component.Base
		timer *timer.Timer
		app   pud.PudApp
		Stats *protos.Stats
	}

	// UserMessage represents a message that user sent
	UserMessage struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	// Stats exports the room status
	Stats struct {
		outboundBytes int
		inboundBytes  int
	}

	// RPCResponse represents a rpc message
	RPCResponse struct {
		Msg string `json:"msg"`
	}

	// SendRPCMsg represents a rpc message
	SendRPCMsg struct {
		ServerID string `json:"serverId"`
		Route    string `json:"route"`
		Msg      string `json:"msg"`
	}

	// NewUser message will be received when new user join room
	NewUser struct {
		Content string `json:"content"`
	}

	// AllMembers contains all members uid
	AllMembers struct {
		Members []string `json:"members"`
	}

	// JoinResponse represents the result of joining room
	JoinResponse struct {
		Code   int    `json:"code"`
		Result string `json:"result"`
	}
)

// NewRoom returns a new room
func NewRoom(app pud.PudApp) *Room {
	return &Room{
		app:   app,
		Stats: &protos.Stats{},
	}
}

// Init runs on service initialization
func (r *Room) Init() {
	r.app.GroupCreate(context.Background(), "room")
}

// AfterInit component lifetime callback
func (r *Room) AfterInit() {
	r.timer = pud.NewTimer(time.Minute, func() {
		count, _ := r.app.GroupCountMembers(context.Background(), "room")
		log.Infof("room user count %v", count)
	})
}

// Entry is the entrypoint
func (r *Room) Entry(ctx context.Context, msg []byte) (*protos.JoinResponse, error) {
	s := r.app.GetSessionFromCtx(ctx)

	err := s.Bind(ctx, uuid.New().String())
	if err != nil {
		log.Errorf("Failed to bind session %v", err.Error())
		return nil, pud.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	return &protos.JoinResponse{Result: "ok"}, nil
}

// GetSessionData gets the session data
func (r *Room) GetSessionData(ctx context.Context) (*SessionData, error) {
	s := r.app.GetSessionFromCtx(ctx)
	return &SessionData{
		Data: s.GetData(),
	}, nil
}

// SetSessionData sets the session data
func (r *Room) SetSessionData(ctx context.Context, data *SessionData) ([]byte, error) {
	s := r.app.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		log.Errorf("Failed to set session data %v", err)
		return nil, err
	}
	err = s.PushToFront(ctx)
	if err != nil {
		return nil, err
	}
	return []byte("success"), nil
}

// Join room
func (r *Room) Join(ctx context.Context) (*protos.JoinResponse, error) {
	s := r.app.GetSessionFromCtx(ctx)
	err := r.app.GroupAddMember(ctx, "room", s.UID())
	if err != nil {
		log.Errorf("Failed to join room %v", err)
		return nil, err
	}
	members, err := r.app.GroupMembers(ctx, "room")
	if err != nil {
		log.Errorf("Failed to get members %v", err)
		return nil, err
	}
	s.Push("onMembers", &protos.AllMembers{Members: members})
	err = r.app.GroupBroadcast(ctx, "connector", "room", "onNewUser", &protos.NewUser{Content: fmt.Sprintf("New user: %d", s.ID())})
	if err != nil {
		log.Errorf("Failed to broadcast onNewUser %v", err)
		return nil, err
	}
	return &protos.JoinResponse{Result: "success"}, nil
}

// Message sync last message to all members
func (r *Room) Message(ctx context.Context, msg *protos.UserMessage) {
	err := r.app.GroupBroadcast(ctx, "connector", "room", "onMessage", msg)
	if err != nil {
		log.Errorf("Error broadcasting message %v", err)
	}
}

// SendRPC sends rpc
func (r *Room) SendRPC(ctx context.Context, msg *protos.SendRPCMsg) (*protos.RPCRes, error) {
	ret := &protos.RPCRes{}
	err := r.app.RPCTo(ctx, msg.ServerId, msg.Route, ret, &protos.RPCMsg{Msg: msg.Msg})
	if err != nil {
		log.Errorf("Failed to execute RPCTo %v - %v %v", msg.ServerId, msg.Route, err)
		return nil, pud.Error(err, "RPC-000")
	}
	return ret, nil
}

// MessageRemote just echoes the given message
func (r *Room) MessageRemote(ctx context.Context, msg *protos.UserMessage, b bool, s string) (*protos.UserMessage, error) {
	return msg, nil
}
