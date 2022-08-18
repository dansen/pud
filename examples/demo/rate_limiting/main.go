package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/dansen/pud"
	"github.com/dansen/pud/acceptor"
	"github.com/dansen/pud/acceptorwrapper"
	"github.com/dansen/pud/component"
	"github.com/dansen/pud/config"
	"github.com/dansen/pud/examples/demo/rate_limiting/services"
	"github.com/dansen/pud/metrics"
)

func createAcceptor(port int, reporters []metrics.Reporter) acceptor.Acceptor {

	// 5 requests in 1 minute. Doesn't make sense, just to test
	// rate limiting
	vConfig := viper.New()
	vConfig.Set("pitaya.conn.ratelimiting.limit", 5)
	vConfig.Set("pitaya.conn.ratelimiting.interval", time.Minute)
	pConfig := config.NewConfig(vConfig)

	rateLimitConfig := config.NewRateLimitingConfig(pConfig)

	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", port))
	return acceptorwrapper.WithWrappers(
		tcp,
		acceptorwrapper.NewRateLimitingWrapper(reporters, *rateLimitConfig))
}

var app pud.Pitaya

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := "room"

	flag.Parse()

	config := config.NewDefaultBuilderConfig()
	builder := pud.NewDefaultBuilder(true, svType, pud.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(createAcceptor(*port, builder.MetricsReporters))

	app = builder.Build()

	defer app.Shutdown()

	room := services.NewRoom()
	app.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.Start()
}
