package console_test

import (
	"bytes"
	"errors"
	"testing"

	"kernel/internal/console"
)

// Step 01 -- the kernel console. Run:  go test ./internal/console -run Step01 -v

func TestStep01_Writef_WritesTheFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	c := console.New(&buf, console.Info)

	if err := c.Writef(console.Info, "boot %s", "kernel"); err != nil {
		t.Fatalf("Writef: unexpected error: %v", err)
	}

	if got := buf.String(); got != "boot kernel\n" {
		t.Fatalf("Writef output = %q, want %q", got, "boot kernel\n")
	}
}

func TestStep01_Write_WritesTheMessageAsAOneLine(t *testing.T) {
	var buf bytes.Buffer
	c := console.New(&buf, console.Info)

	if err := c.Write(console.Info, "hello"); err != nil {
		t.Fatalf("Write: unexpected error: %v", err)
	}

	if got := buf.String(); got != "hello\n" {
		t.Fatalf("Write output = %q, want %q", got, "hello\n")
	}
}

func TestStep01_Console_DropsLinesBelowTheMinimumLevel(t *testing.T) {
	var buf bytes.Buffer
	c := console.New(&buf, console.Warn)

	if err := c.Writef(console.Debug, "hidden debug line"); err != nil {
		t.Fatalf("Writef(Debug): unexpected error: %v", err)
	}
	if err := c.Writef(console.Info, "hidden info line"); err != nil {
		t.Fatalf("Writef(Info): unexpected error: %v", err)
	}
	if err := c.Writef(console.Warn, "shown warn line"); err != nil {
		t.Fatalf("Writef(Warn): unexpected error: %v", err)
	}

	if got := buf.String(); got != "shown warn line\n" {
		t.Fatalf("Writef output = %q, want only the warn line %q", got, "shown warn line\n")
	}
}

func TestStep01_Console_WritesErrorLinesEvenWhenTheMinimumIsHigher(t *testing.T) {
	var buf bytes.Buffer
	c := console.New(&buf, console.Info)

	if err := c.Write(console.Error, "oops"); err != nil {
		t.Fatalf("Write(Error): unexpected error: %v", err)
	}

	if got := buf.String(); got != "oops\n" {
		t.Fatalf("Write output = %q, want %q", got, "oops\n")
	}
}

func TestStep01_Console_ReturnsTheSinkError(t *testing.T) {
	c := console.New(errWriter{}, console.Info)

	if err := c.Writef(console.Info, "hello"); err == nil {
		t.Fatal("Writef: expected the sink error to propagate, got nil")
	}
}

// errWriter is a sink that always fails, to check error propagation.
type errWriter struct{}

func (errWriter) Write(_ []byte) (int, error) { return 0, errors.New("sink: boom") }
