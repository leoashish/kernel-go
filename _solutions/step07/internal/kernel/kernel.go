// REFERENCE SOLUTION for Step 07. Try it yourself first!
// Package kernel wires the miniature kernel together: console, syscall table, scheduler, processes and context.
package kernel

import (
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

// Run performs one timer interrupt: the timer fires, the scheduler runs the next process, and the CPU runs it —
// here, running the boot process means writing the boot line through the write syscall. The timer's tick handler IS
// the scheduler's run-next, so firing the timer performs the run.
func (k *Kernel) Run() error {
	if err := k.timer.Fire(); err != nil {
		return err
	}
	return k.syscalls.Dispatch(syscall.SYS_WRITE, syscall.StdoutFd, BootLine)
}
