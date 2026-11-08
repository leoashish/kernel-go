// Package kernel wires the miniature kernel together: console, syscall table, scheduler, processes and context.
package kernel

import (
	"errors"
	"os"
	"time"

	"kernel/internal/console"
	"kernel/internal/process"
	"kernel/internal/scheduler"
	"kernel/internal/syscall"
)

// BootLine is the line the initial user process writes through the write syscall when the kernel boots. GIVEN.
const BootLine = "hello from userland"

// Kernel is the miniature kernel: console + syscall table + scheduler + processes wired together.
// GIVEN.
type Kernel struct {
	console   *console.Console
	syscalls  *syscall.DispatchTable
	scheduler *scheduler.Scheduler
	procs     *process.ProcessTable
	timer     *scheduler.Timer
	entry     int
}

// New returns a freshly wired miniature kernel. GIVEN.
func New() *Kernel {
	tm := scheduler.New(time.Millisecond, func() error { return nil })
	procs := process.New()
	sched := scheduler.NewScheduler(tm, procs)
	return &Kernel{
		console:   console.New(os.Stdout, console.Info),
		syscalls:  syscall.New(),
		scheduler: sched,
		procs:     procs,
		timer:     tm,
		entry:     0,
	}
}

// Boot registers the built-in syscalls, spawns the initial user process and runs it through the timer interrupt.
func (k *Kernel) Boot() error {
	k.syscalls.Register(syscall.SYS_WRITE, syscall.Write)
	k.syscalls.Register(syscall.SYS_EXIT, syscall.Exit)
	if _, err := k.procs.Create(1, process.Registers{PC: k.entry}); err != nil {
		return err
	}
	if err := k.scheduler.Add(1); err != nil {
		return err
	}
	return k.Run()
}

// Run performs one timer interrupt: the scheduler runs the next process, and the kernel executes that process —
// here, running the boot process means writing the boot line through the write syscall. If the executed action asks
// the process to exit, the kernel tears the process down (marks it exited).
func (k *Kernel) Run() error {
	proc, err := k.scheduler.RunNext()
	if err != nil {
		return err
	}
	return k.executeBoot(proc)
}

// executeBoot runs the action of the booted process via the syscall table; if the action ends in an exit, tear the
// process down.
func (k *Kernel) executeBoot(proc *process.Process) error {
	err := k.syscalls.Dispatch(syscall.SYS_WRITE, syscall.StdoutFd, BootLine)
	var exitErr *syscall.ExitError
	if errors.As(err, &exitErr) {
		// teardown: the process exited; mark it exited so the scheduler never runs it again.
		if _, perr := k.procs.Exit(proc.Pid); perr != nil {
			return perr
		}
		return nil
	}
	return err
}
