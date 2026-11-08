// Package cpu implements context switches: save/restore of the CPU context between processes.
package cpu

import (
	"kernel/internal/process"
)

// Context is one saved execution context: which process is running and its registers.
// GIVEN.
type Context struct {
	Process *process.Process
	Regs    process.Registers
}

// The context the CPU is currently running. GIVEN (the miniature keeps it in the cpu package for this step).
var current Context

// Save captures the context currently in the CPU and returns it. GIVEN.
func Save() Context {
	current = Context{Process: current.Process, Regs: current.Regs}
	return current
}

// Restore loads a captured context into the CPU and returns it as the current context. GIVEN.
func Restore(c Context) Context {
	current = c
	return current
}

// Switch loads context c into the CPU, saving the context that was running before; it returns that previous context
// so the caller can Restore it when switching back. GIVEN.
func Switch(c Context) Context {
	prev := current
	current = c
	return prev
}
