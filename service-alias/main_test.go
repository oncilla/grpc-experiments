package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var buf bytes.Buffer
	out = &buf

	server := newServerCmd()
	go func() {
		server.Execute()
	}()

	time.Sleep(time.Second)

	fmt.Fprintln(out, "\n-- Regular Run:")
	client := newClientCmd()
	if err := client.Execute(); err != nil {
		t.Fatal(err)
	}
	fmt.Fprintln(out, "\n-- Legacy Run:")
	client.SetArgs([]string{"--legacy"})
	if err := client.Execute(); err != nil {
		t.Fatal(err)
	}

	expected := `
-- Regular Run:
Server: ping
Response: ping

-- Legacy Run:
Sending legacy request
Server: legacy
Response: legacy
`
	if buf.String() != expected {
		t.Fatal(buf.String())
	}
}
