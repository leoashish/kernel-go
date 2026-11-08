// Package console implements the kernel's output device: the single sink for boot lines, driver messages and oops
// reports. Step 01 gives the kernel a console; later steps (syscalls, interrupts, drivers) just call Writef on it.
package console

import (
	"fmt"
	"io"
)

// Level is the severity of a kernel log line. Debug(0) < Info(1) < Warn(2) < Error(3); a Console emits a line when
// its level is greater than or equal to the Console's minimum level.
type Level int

// The kernel's log levels, in increasing order of severity.
const (
	Debug Level = iota
	Info
	Warn
	Error
)

// Console writes kernel log lines to an underlying sink, honouring a minimum level.
//
// STEP 01: the whole device. It is a thin level-gated writer: Write/Writef format-and-emit through the sink, dropping
// lines below the minimum level. Later steps (syscalls, interrupts, drivers) just call Writef and never talk to the
// sink themselves.
type Console struct {
	out      io.Writer
	minLevel Level
}

// New returns a Console that writes to out and emits every line at or above minLevel.
func New(out io.Writer, minLevel Level) *Console {
	return &Console{out: out, minLevel: minLevel}
}

// Write emits one line at the given level, followed by a newline. Lines below the Console's minimum level are
// dropped (a no-op returning nil). The error of the underlying sink is returned to the caller.
func (c *Console) Write(level Level, msg string) error {
	if level < c.minLevel {
		return nil
	}
	_, err := fmt.Fprintln(c.out, msg)
	return err
}

// Writef emits one line built with fmt.Sprintf(format, args...) followed by a newline, at the given level, honouring
// the minimum level. The error of the underlying sink is returned to the caller.
func (c *Console) Writef(level Level, format string, args ...any) error {
	if level < c.minLevel {
		return nil
	}
	_, err := fmt.Fprintf(c.out, format+"\n", args...)
	return err
}
