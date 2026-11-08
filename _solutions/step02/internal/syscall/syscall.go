package syscall

import (
	"errors"
	"fmt"
	"os"

	"kernel/internal/console"
)

// REFERENCE SOLUTION for Step 02. Try it yourself first!

// ErrNoSys: the syscall number is not known to this kernel (the Linux analogue: errno ENOSYS). GIVEN.
var ErrNoSys = errors.New("syscall: ENOSYS: unknown syscall")

// ErrBadFd: the file descriptor is not open in this process (the Linux analogue: errno EBADF). GIVEN.
var ErrBadFd = errors.New("syscall: EBADF: bad file descriptor")

// ExitError signals that a process asked to terminate with the given code. In the real kernel the exit syscall
// never returns; in this miniature the handler reports the request as an error. GIVEN.
type ExitError struct {
	Code int
}

// Error implements error.
func (e *ExitError) Error() string { return fmt.Sprintf("syscall: exit %d", e.Code) }

// StdoutFd / StderrFd are the file descriptors user programs write their output to.
const (
	StdoutFd = 1
	StderrFd = 2
)

// SYS_WRITE is the syscall number of write (the classic first syscall every process runs). GIVEN.
const SYS_WRITE = 1

// SYS_EXIT is the syscall number of exit (the syscall that ends a process). GIVEN.
const SYS_EXIT = 60

// Handler handles one syscall invocation. args are the arguments after the syscall number. In this miniature the
// ABI is untyped: numeric syscalls pass their numbers as ints, and the console-writing SYS_WRITE passes its line as
// a string next to the fd. GIVEN.
type Handler func(args ...any) error

// DispatchTable maps syscall numbers to handlers.
type DispatchTable struct {
	handlers map[int]Handler
}

// New returns an empty dispatch table.
func New() *DispatchTable {
	return &DispatchTable{handlers: map[int]Handler{}}
}

// Register installs h for syscall number num (replacing any previous handler).
func (t *DispatchTable) Register(num int, h Handler) {
	t.handlers[num] = h
}

// Lookup returns the handler registered for num, and whether the number is known.
func (t *DispatchTable) Lookup(num int) (Handler, bool) {
	h, ok := t.handlers[num]
	return h, ok
}

// Dispatch routes a syscall: it looks the number up and invokes the handler with the args. An unknown number yields
// ErrNoSys (the kernel does not know that syscall yet).
func (t *DispatchTable) Dispatch(num int, args ...any) error {
	h, ok := t.Lookup(num)
	if !ok {
		return ErrNoSys
	}
	return h(args...)
}

// Write handles SYS_WRITE: it writes one console line to the console sink on the given fd (1 = stdout, 2 = stderr;
// in this miniature both map to the console). The kernel console is the only sink today; later steps add real file
// descriptors.
func Write(args ...any) error {
	if len(args) != 2 {
		return fmt.Errorf("syscall: write: expected (fd, line), got %d argument(s)", len(args))
	}
	fd, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("syscall: write: fd must be an int, got %T", args[0])
	}
	line, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("syscall: write: line must be a string, got %T", args[1])
	}
	if fd != StdoutFd && fd != StderrFd {
		return ErrBadFd
	}
	return console.New(os.Stdout, console.Info).Writef(console.Info, "%s", line)
}

// Exit handles SYS_EXIT: it reports the process termination with the given code. The real kernel never returns from
// this syscall; the handler models the request as an error (process teardown comes in Step 03).
func Exit(args ...any) error {
	if len(args) < 1 {
		return fmt.Errorf("syscall: exit: expected the exit code, got none")
	}
	code, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("syscall: exit: code must be an int, got %T", args[0])
	}
	return &ExitError{Code: code}
}
