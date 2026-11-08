// REFERENCE SOLUTION for Step 04. Try it yourself first!
// Package scheduler implements the kernel's timer interrupt and run queue.
package scheduler

import (
	"errors"
	"sync"
	"time"

	"kernel/internal/process"
)

// ErrEmpty: the run queue has no process to run (spurious or exhausted interrupt). GIVEN.
var ErrEmpty = errors.New("scheduler: empty run queue")

// Timer models a periodic timer interrupt source: the kernel arms a timer, and when the interrupt fires it invokes
// the tick handler (the scheduler's interrupt handler). GIVEN.
type Timer struct {
	mu       sync.Mutex
	interval time.Duration
	tick     func() error
	stopped  bool
}

// New returns a Timer that fires every interval and calls tick.
func New(interval time.Duration, tick func() error) *Timer {
	return &Timer{interval: interval, tick: tick}
}

// Interval returns the timer's period.
func (t *Timer) Interval() time.Duration { return t.interval }

// Fire delivers the timer interrupt: it invokes the tick handler unless the timer is stopped.
func (t *Timer) Fire() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped {
		return nil
	}
	return t.tick()
}

// Stop stops the timer: later Fires become no-ops.
func (t *Timer) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopped = true
	return nil
}

// Stopped reports whether the timer has been stopped. GIVEN.
func (t *Timer) Stopped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.stopped
}

// Scheduler holds the run queue of pids and runs processes round-robin on timer interrupts.
type Scheduler struct {
	mu    sync.Mutex
	queue []int
	timer *Timer
	procs *process.ProcessTable
}

// NewScheduler returns a Scheduler whose timer tick runs the next process in the queue (FIFO).
func NewScheduler(t *Timer, procs *process.ProcessTable) *Scheduler {
	s := &Scheduler{queue: []int{}, timer: t, procs: procs}
	t.tick = func() error { _, err := s.RunNext(); return err }
	return s
}

// Add enqueues a pid to run when the timer fires.
func (s *Scheduler) Add(pid int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, pid)
	return nil
}

// Len returns the number of processes currently queued.
func (s *Scheduler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue)
}

// RunNext is the timer's tick handler: it pops the front of the queue and returns the process to run (the CPU
// dispatch happens in the kernel layer).
func (s *Scheduler) RunNext() (*process.Process, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		return nil, ErrEmpty
	}
	pid := s.queue[0]
	s.queue = s.queue[1:]
	return s.procs.Get(pid)
}
