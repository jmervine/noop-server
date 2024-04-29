package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/net/server"
)

var cfg *config.Config

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	var err error
	cfg, err = config.Init(os.Args, Version)

	if err != nil {
		log.Fatal(err)
	}

	log.SetPrefix(fmt.Sprintf("app=%s ", cfg.App))
	log.Printf("on=startup %s\n", cfg.ToString())
	log.Fatal(server.Start(cfg))
}
