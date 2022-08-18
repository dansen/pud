package services

import (
	"context"
	"github.com/dansen/pud/defaultlog/log"

	"github.com/dansen/pud/component"
	"github.com/dansen/pud/examples/demo/worker/protos"
)

// Metagame server
type Metagame struct {
	component.Base
}

// LogRemote logs argument when called
func (m *Metagame) LogRemote(ctx context.Context, arg *protos.Arg) (*protos.Response, error) {
	log.Infof("argument %+v\n", arg)
	return &protos.Response{Code: 200, Msg: "ok"}, nil
}
