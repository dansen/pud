package services

import (
	"context"

	"github.com/dansen/pud/logger/log"

	"github.com/dansen/pud"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/examples/demo/worker/protos"
)

// Room server
type Room struct {
	component.Base
	app pud.PudApp
}

// NewRoom ctor
func NewRoom(app pud.PudApp) *Room {
	return &Room{app: app}
}

// CallLog makes ReliableRPC to metagame LogRemote
func (r *Room) CallLog(ctx context.Context, arg *protos.Arg) (*protos.Response, error) {
	route := "metagame.metagame.logremote"
	reply := &protos.Response{}
	jid, err := r.app.ReliableRPC(route, nil, reply, arg)
	if err != nil {
		log.Infof("failed to enqueue rpc: %q", err)
		return nil, err
	}

	log.Infof("enqueue rpc job: %d", jid)
	return &protos.Response{Code: 200, Msg: "ok"}, nil
}
