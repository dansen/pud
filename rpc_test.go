// Copyright (c) TFG Co. All Rights Reserved.
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

package pud

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/dansen/pud/cluster"
	clustermocks "github.com/dansen/pud/cluster/mocks"
	"github.com/dansen/pud/config"
	"github.com/dansen/pud/conn/codec"
	"github.com/dansen/pud/conn/message"
	"github.com/dansen/pud/constants"
	"github.com/dansen/pud/pipeline"
	"github.com/dansen/pud/protos"
	"github.com/dansen/pud/protos/test"
	"github.com/dansen/pud/route"
	"github.com/dansen/pud/router"
	serializemocks "github.com/dansen/pud/serialize/mocks"
	"github.com/dansen/pud/service"
	sessionmocks "github.com/dansen/pud/session/mocks"
)

func TestDoSendRPCNotInitialized(t *testing.T) {
	config := config.NewDefaultBuilderConfig()
	app := NewDefaultApp(true, "testtype", Standalone, map[string]string{}, *config).(*App)
	err := app.doSendRPC(nil, "", "", nil, nil)
	assert.Equal(t, constants.ErrRPCServerNotInitialized, err)
}

func TestDoSendRPC(t *testing.T) {
	config := config.NewDefaultBuilderConfig()
	app := NewDefaultApp(true, "testtype", Cluster, map[string]string{}, *config).(*App)
	app.server.ID = "myserver"
	app.rpcServer = &cluster.NatsRPCServer{}
	tables := []struct {
		name     string
		routeStr string
		reply    proto.Message
		arg      proto.Message
		err      error
	}{
		{"bad_route", "badroute", &test.SomeStruct{}, nil, route.ErrInvalidRoute},
		{"no_server_type", "bla.bla", &test.SomeStruct{}, nil, constants.ErrNoServerTypeChosenForRPC},
		{"nonsense_rpc", "mytype.bla.bla", &test.SomeStruct{}, nil, constants.ErrNonsenseRPC},
		{"success", "bla.bla.bla", &test.SomeStruct{}, &test.SomeStruct{A: 1}, nil},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			ctx := context.Background()
			if table.err == nil {
				packetEncoder := codec.NewPomeloPacketEncoder()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockSerializer := serializemocks.NewMockSerializer(ctrl)
				mockSD := clustermocks.NewMockServiceDiscovery(ctrl)
				mockRPCClient := clustermocks.NewMockRPCClient(ctrl)
				mockRPCServer := clustermocks.NewMockRPCServer(ctrl)
				messageEncoder := message.NewMessagesEncoder(false)
				sessionPool := sessionmocks.NewMockSessionPool(ctrl)
				router := router.New()
				handlerPool := service.NewHandlerPool()
				svc := service.NewRemoteService(mockRPCClient, mockRPCServer, mockSD, packetEncoder, mockSerializer, router, messageEncoder, &cluster.Server{}, sessionPool, pipeline.NewHandlerHooks(), handlerPool)
				assert.NotNil(t, svc)
				app.remoteService = svc
				app.server.ID = "notmyserver"
				b, err := proto.Marshal(&test.SomeStruct{A: 1})
				assert.NoError(t, err)
				mockSD.EXPECT().GetServer("myserver").Return(&cluster.Server{}, nil)
				mockRPCClient.EXPECT().Call(ctx, protos.RPCType_User, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&protos.Response{Data: b}, nil)
			}
			err := app.RPCTo(ctx, "myserver", table.routeStr, table.reply, table.arg)
			assert.Equal(t, table.err, err)
		})
	}
}
