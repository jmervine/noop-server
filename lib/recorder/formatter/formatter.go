package formatter

import "github.com/jmervine/noop-server/lib/recorder/formatter/noop_client"

type Default struct {
	// Change this to change the default formatter.
	noop_client.NoopClient
}
