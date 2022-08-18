package services

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/dansen/pud"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/examples/demo/worker/protos"
)

// Worker server
type Worker struct {
	component.Base
}

// Configure starts workers and register rpc job
func (w *Worker) Configure(app pud.PudApp) {
	app.StartWorker()
	app.RegisterRPCJob(&RPCJob{app: app})
}

// RPCJob implements worker.RPCJob
type RPCJob struct {
	app pud.PudApp
}

// ServerDiscovery returns a serverID="", meaning any server
// is ok
func (r *RPCJob) ServerDiscovery(
	route string,
	rpcMetadata map[string]interface{},
) (serverID string, err error) {
	return "", nil
}

// RPC calls pitaya's rpc
func (r *RPCJob) RPC(
	ctx context.Context,
	serverID, routeStr string,
	reply, arg proto.Message,
) error {
	return r.app.RPCTo(ctx, serverID, routeStr, reply, arg)
}

// GetArgReply returns reply and arg of LogRemote,
// since we have no other methods in this example
func (r *RPCJob) GetArgReply(
	route string,
) (arg, reply proto.Message, err error) {
	return &protos.Arg{}, &protos.Response{}, nil
}
