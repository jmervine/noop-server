package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const DEFAULT_NAME = "noop-server"
const DEFAULT_RECORD_TARGET = "record.txt"
const DEFAULT_RECORD_TARGET_APPEND = ".bak"

type Config struct {
	App  string
	Port string
	Addr string

	CertPrivatePath string
	CertKeyPath     string
	CertCAPath      string

	Verbose bool

	Record       bool
	RecordTarget string
}

func Init(args []string) *Config {
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
			Value:       "localhost",
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
		// TODO: Support record formats: noop-client,csv,json,yaml
		&cli.BoolFlag{
			Name:        "record",
			Aliases:     []string{"r"},
			Usage:       "Record results to a file",
			Value:       false,
			Destination: &c.Record,
		},
		&cli.StringFlag{
			Name:        "record-target",
			Aliases:     []string{"t"},
			Usage:       "Record results to a file",
			Value:       DEFAULT_RECORD_TARGET,
			Destination: &c.RecordTarget,
		},
	}
	app.Action = func(_ *cli.Context) error {
		return nil
	}

	if err := app.Run(args); err != nil {
		log.Fatal(err)
	}

	return &c
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

// Create RecordTarget file, if Record is configured. If it already
// exists, back up the old one.
//
// This returned 'os.File' must be closed when you're done with it.
func (c Config) RecordFile() (*os.File, error) {
	if !c.Record {
		return nil, nil
	}

	// Workflow
	// 1. Does file exist
	// 2. IF YES, backup and create
	// 3. IF NO, create
	var err error

	// Check to see if the file exists
	_, err = os.Stat(c.RecordTarget)
	if err == nil {
		// The file exists, so back it up.
		bak := c.RecordTarget + DEFAULT_RECORD_TARGET_APPEND
		err = createBackup(c.RecordTarget, bak)
		if err != nil {
			return nil, err
		}
	} else {
		// Return errors that aren't existance check errors
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	// Assume the file doesn't exist and create and return it.
	return os.Create(c.RecordTarget)
}

// REF: https://stackoverflow.com/a/62179184
// This moves the file, removing the original file.
func createBackup(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}
