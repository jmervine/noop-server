package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jmervine/noop-server/lib/records/formatter"
)

func baseConfig() *Config {
	return &Config{
		App:  "test-noop-server",
		Port: "3333",
		Addr: "localhost",
	}
}

func TestConfig_validate(t *testing.T) {
	config := baseConfig()

	err := config.validate()
	if err != nil {
		t.Errorf("expected %#v to be valid", config)
	}

	config.recordFormat = "ack"

	err = config.validate()
	if err == nil {
		t.Errorf("expected %#v to be invalid", config)
	}

	config = baseConfig()
	config.Record = true
	config.StreamRecord = true

	err = config.validate()
	if err == nil {
		t.Errorf("expected %#v to be invalid", config)
	}
}

func TestConfig_RecordFormater(t *testing.T) {
	config := baseConfig()

	d := &formatter.Default{}
	if !reflect.DeepEqual(config.RecordFormatter(), d) {
		t.Errorf("Expected RecordFormatter to return Default, got %#v", config.RecordFormatter())
	}

	config.recordFormat = "json"
	j := &formatter.Json{}
	if !reflect.DeepEqual(config.RecordFormatter(), j) {
		t.Errorf("Expected RecordFormatter to return %#v, got %#v", j, config.RecordFormatter())
	}
}

func TestConfig_MaxProcs(t *testing.T) {
	config := baseConfig()

	if config.MaxProcs() != 1 {
		t.Error("Expected MaxProcs to be 1, got", config.MaxProcs())
	}
}

func TestConfig_Listener(t *testing.T) {
	config := baseConfig()

	if config.Listener() != "localhost:3333" {
		t.Error("Expected MaxProcs to be localhost:3333, got", config.Listener())
	}
}

func TestConfig_TLSEnabled(t *testing.T) {
	config := baseConfig()

	if config.TLSEnabled() {
		t.Error("Expected TLSEnabled to be false, was true")
	}

	config.CertKeyPath = "somefile"
	config.CertPrivatePath = "somefile"

	if !config.TLSEnabled() {
		t.Error("Expected TLSEnabled to be true, was false")
	}
}

func TestConfig_MTLSEnabled(t *testing.T) {
	config := baseConfig()
	config.CertKeyPath = "somefile"
	config.CertPrivatePath = "somefile"

	if config.MTLSEnabled() {
		t.Error("Expected MTLSEnabled to be false, was true")
	}

	config.CertCAPath = "somefile"

	if !config.MTLSEnabled() {
		t.Error("Expected MTLSEnabled to be true, was false")
	}
}

func TestConfig_Recording(t *testing.T) {
	config := baseConfig()
	if config.Recording() {
		t.Error("Expected Recording to be false, but was true")
	}

	config.Record = true
	if !config.Recording() {
		t.Error("Expected Recording to be true, but was false")
	}
}

func TestConfig_ToString(t *testing.T) {
	config := baseConfig()
	expect := "addr=localhost port=3333 mtls=false ssl=false verbose=false record=false"

	got := config.ToString()
	if !strings.Contains(got, expect) {
		t.Errorf("Expected '%s' to contain '%s'", got, expect)
	}

	config.Record = true
	got = config.ToString()
	if !strings.Contains(got, "record-target") {
		t.Errorf("Expected '%s' to contain 'record-target'", got)
	}
	if !strings.Contains(got, "record-format") {
		t.Errorf("Expected '%s' to contain 'record-format'", got)
	}
}
