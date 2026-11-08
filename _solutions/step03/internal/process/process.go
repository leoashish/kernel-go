package process

import (
	"errors"
	"sync"
)

// REFERENCE SOLUTION for Step 03. Try it yourself first!

// ErrNoSuchProcess: the pid does not belong to a live process. GIVEN.
var ErrNoSuchProcess = errors.New("process: no such process")

// ErrProcessExists: a process already runs with the pid we were asked to use. GIVEN.
var ErrProcessExists = errors.New("process: already exists")

// NumGprs is how many general-purpose registers a simulated process has. GIVEN.
const NumGprs = 16

// State is what a process is doing. GIVEN.
type State int

// The kernel's process states, in increasing order of progress. GIVEN.
const (
	Ready State = iota
	Running
	Exited
)

// Registers holds the CPU state of a process: a program counter plus general-purpose registers. GIVEN.
type Registers struct {
	PC  int
	GPR [NumGprs]int
}

// Process is one task in the kernel. GIVEN.
type Process struct {
	Pid   int
	State State
	Regs  Registers
}

// ProcessTable holds the live processes of the kernel. GIVEN.
type ProcessTable struct {
	mu      sync.Mutex
	procs   map[int]*Process
	nextPid int
}

// New returns an empty process table.
func New() *ProcessTable {
	return &ProcessTable{procs: map[int]*Process{}, nextPid: 1}
}

// Create makes a process with the given pid and registers.
func (t *ProcessTable) Create(pid int, regs Registers) (*Process, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.procs[pid]; ok {
		return nil, ErrProcessExists
	}
	p := &Process{Pid: pid, State: Ready, Regs: regs}
	t.procs[pid] = p
	return p, nil
}

// Fork spawns a child with the given pid, copying the registers of parent; the child starts Ready.
func (t *ProcessTable) Fork(pid int, parent *Process) (*Process, error) {
	if parent == nil {
		return nil, ErrNoSuchProcess
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.procs[parent.Pid]; !ok {
		return nil, ErrNoSuchProcess
	}
	if _, ok := t.procs[pid]; ok {
		return nil, ErrProcessExists
	}
	child := &Process{Pid: pid, State: Ready, Regs: parent.Regs}
	t.procs[pid] = child
	return child, nil
}

// Get returns the process for the given pid, or an error satisfying errors.Is(err, ErrNoSuchProcess).
func (t *ProcessTable) Get(pid int) (*Process, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.procs[pid]
	if !ok {
		return nil, ErrNoSuchProcess
	}
	return p, nil
}

// Exec replaces the address space of the process with the entry point: the program counter moves to entry and the
// process goes back to Ready (waiting for the scheduler).
func (t *ProcessTable) Exec(pid int, entry int) (*Process, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.procs[pid]
	if !ok {
		return nil, ErrNoSuchProcess
	}
	p.Regs.PC = entry
	p.State = Ready
	return p, nil
}
