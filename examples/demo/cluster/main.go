package main

import (
	"context"
	"flag"
	"fmt"

	"strings"

	"github.com/dansen/pud"
	"github.com/dansen/pud/acceptor"
	"github.com/dansen/pud/cluster"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/config"
	"github.com/dansen/pud/examples/demo/cluster/services"
	"github.com/dansen/pud/groups"
	"github.com/dansen/pud/internal/generic/log"
	"github.com/dansen/pud/internal/generic/log/lowlevel"
	"github.com/dansen/pud/route"
)

var app pud.Pitaya

func configureBackend() {
	room := services.NewRoom(app)
	app.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.RegisterRemote(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)
}

func configureFrontend(port int) {
	app.Register(services.NewConnector(app),
		component.WithName("connector"),
		component.WithNameFunc(strings.ToLower),
	)

	app.RegisterRemote(services.NewConnectorRemote(app),
		component.WithName("connectorremote"),
		component.WithNameFunc(strings.ToLower),
	)

	err := app.AddRoute("room", func(
		ctx context.Context,
		route *route.Route,
		payload []byte,
		servers map[string]*cluster.Server,
	) (*cluster.Server, error) {
		// will return the first server
		for k := range servers {
			return servers[k], nil
		}
		return nil, nil
	})

	if err != nil {
		fmt.Printf("error adding route %s\n", err.Error())
	}

	err = app.SetDictionary(map[string]uint16{
		"connector.getsessiondata": 1,
		"connector.setsessiondata": 2,
		"room.room.getsessiondata": 3,
		"onMessage":                4,
		"onMembers":                5,
	})

	if err != nil {
		fmt.Printf("error setting route dictionary %s\n", err.Error())
	}
}

func main() {
	log.SetLogger(lowlevel.NewDefaultLogger())
	port := flag.Int("port", 3250, "the port to listen")
	svType := flag.String("type", "connector", "the server type")
	isFrontend := flag.Bool("frontend", true, "if server is frontend")

	flag.Parse()

	builderConfig := config.NewDefaultBuilderConfig()
	builder := pud.NewDefaultBuilder(*isFrontend, *svType, pud.Cluster, map[string]string{}, *builderConfig)

	if *isFrontend {
		// 前端接受客户端
		tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port))
		builder.AddAcceptor(tcp)
	}
	builder.Groups = groups.NewMemoryGroupService(*config.NewDefaultMemoryGroupConfig())
	app = builder.Build()

	//TODO: Oelze pitaya.SetSerializer(protobuf.NewSerializer())

	defer app.Shutdown()

	if !*isFrontend {
		configureBackend()
	} else {
		configureFrontend(*port)
	}

	app.Start()
}
