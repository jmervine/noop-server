package responder

import (
	"testing"
)

func TestResponders_Load(t *testing.T) {
	// Root for this test is 'lib/responder'
	responders, err := Load("../../test/script.yaml")

	if err != nil {
		t.Errorf("Unexpected error: %+v\n", err)
	}

	responder, ok := (*responders)["/endpoint*"]
	if !ok {
		t.Errorf("Expected '/endpoint*', got %+v", responders)
	}

	if responder.Status != 200 {
		t.Errorf("Expected 200, got %d", responder.Status)
	}

	if responder.Sleep != 2000 {
		t.Errorf("Expected 2000, got %d", responder.Status)
	}

	if responder.Response != "script response" {
		t.Errorf("Expected 'script response', got %s", responder.Response)
	}
}
