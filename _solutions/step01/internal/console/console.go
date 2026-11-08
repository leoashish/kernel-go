package console

import (
	"fmt"
	"io"
)

// REFERENCE SOLUTION for Step 01. Try it yourself first!

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
