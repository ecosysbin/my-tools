package main

import (
	"flag"
	"fmt"
	"os"

	"greet/internal/config"
	"greet/internal/handler"
	"greet/internal/svc"

	"github.com/rs/zerolog"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/greet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	log := zerolog.New(os.Stderr)
	log = log.With().Caller().Logger()
	log.Info().Msg("hello world")
	server.Start()
}
