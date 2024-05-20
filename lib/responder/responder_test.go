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

func responders() Responders {
	r := Responders{}
	r["/test1"] = Responder{}
	r["/test3"] = Responder{}
	r["/test5"] = Responder{}
	r["*"] = Responder{}
	return r
}

func BenchmarkResponders_Match_static(b *testing.B) {
	r := responders()

	for n := 0; n < b.N; n++ {
		r.Match("/test1")
	}
}

func BenchmarkResponders_Match_star(b *testing.B) {
	r := responders()

	for n := 0; n < b.N; n++ {
		r.Match("*")
	}
}

func BenchmarkResponders_Match_wildcard(b *testing.B) {
	r := responders()

	for n := 0; n < b.N; n++ {
		r.Match("/test*")
	}
}
