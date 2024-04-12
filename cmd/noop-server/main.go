package main

import (
	gotls "crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/net/server"
)

var (
	cfg  *config.Config
	cert gotls.Certificate
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	cfg = config.Init(os.Args)

	log.SetPrefix(fmt.Sprintf("app=%s ", cfg.App))
	log.Printf("on=startup %s\n", cfg.ToString())
	log.Fatal(server.Start(cfg))
}
