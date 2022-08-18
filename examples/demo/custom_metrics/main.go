package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/dansen/pud"
	"github.com/dansen/pud/acceptor"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/config"
	"github.com/dansen/pud/examples/demo/custom_metrics/services"
)

var app pud.Pitaya

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := "room"
	isFrontend := true

	flag.Parse()

	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigName("config")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}

	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port))

	conf := config.NewConfig(cfg)
	builder := pud.NewBuilderWithConfigs(isFrontend, svType, pud.Cluster, map[string]string{}, conf)
	builder.AddAcceptor(tcp)
	app = builder.Build()

	defer app.Shutdown()

	app.Register(services.NewRoom(app),
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.Start()
}
