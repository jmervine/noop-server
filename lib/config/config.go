package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const DEFAULT_NAME = "noop-server"

type Config struct {
	App  string
	Port string
	Addr string

	CertPrivatePath string
	CertKeyPath     string
	CertCAPath      string

	Verbose bool
}

func Init(args []string) Config {
	c := Config{}

	helpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, t string, d interface{}) {
		helpPrinter(w, t, d)
		os.Exit(0)
	}

	app := cli.NewApp()

	app.Name = DEFAULT_NAME
	app.Usage = "A simple noop server that accepts everything"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			Usage:       "App name for logger output",
			Value:       DEFAULT_NAME,
			EnvVars:     []string{"APP_NAME"},
			Required:    false,
			Destination: &c.App,
		},
		&cli.StringFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			Usage:       "Listener port",
			Value:       "3000",
			EnvVars:     []string{"PORT"},
			Required:    false,
			Destination: &c.Port,
		},
		&cli.StringFlag{
			Name:        "addr",
			Aliases:     []string{"a"},
			Usage:       "Listener address",
			Value:       "0.0.0.0",
			EnvVars:     []string{"ADDR"},
			Required:    false,
			Destination: &c.Addr,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "Enable verbose logging",
			Value:       false,
			EnvVars:     []string{"VERBOSE"},
			Required:    false,
			Destination: &c.Verbose,
		},

		&cli.StringFlag{
			Name:        "private",
			Usage:       "TLS Cert Private file path",
			Required:    false,
			EnvVars:     []string{"TLS_PRIVATE_PATH"},
			Destination: &c.CertPrivatePath,
		},
		&cli.StringFlag{
			Name:        "key",
			Usage:       "TLS Cert Key file path",
			Required:    false,
			EnvVars:     []string{"TLS_KEY_PATH"},
			Destination: &c.CertKeyPath,
		},
		&cli.StringFlag{
			Name:        "ca",
			Usage:       "mTLS Cert Chain file path",
			Required:    false,
			EnvVars:     []string{"MTLS_CA_CHAIN_PATH"},
			Destination: &c.CertCAPath,
		},
	}
	app.Action = func(_ *cli.Context) error {
		return nil
	}

	if err := app.Run(args); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c Config) Listener() string {
	return fmt.Sprintf("%s:%s", c.Addr, c.Port)
}

func (c Config) TLSEnabled() bool {
	return (c.CertKeyPath != "" && c.CertPrivatePath != "")
}

func (c Config) MTLSEnabled() bool {
	return (c.TLSEnabled() && c.CertCAPath != "")
}

func (c Config) ToString() string {
	return fmt.Sprintf(
		"addr=%s port=%s mtls=%v ssl=%v verbose=%v",
		c.Addr, c.Port, c.MTLSEnabled(), c.TLSEnabled(), c.Verbose)
}
