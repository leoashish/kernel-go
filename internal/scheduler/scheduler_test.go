package scheduler_test

import (
	"errors"
	"testing"
	"time"

	"kernel/internal/process"
	"kernel/internal/scheduler"
)

// Step 04 -- the timer interrupt and the scheduler's run queue. Run: go test ./internal/scheduler -run Step04 -v

func TestStep04_Timer_FiresTheTick(t *testing.T) {
	called := false
	tm := scheduler.New(time.Millisecond, func() error {
		called = true
		return nil
	})

	if err := tm.Fire(); err != nil {
		t.Fatalf("Fire: unexpected error: %v", err)
	}
	if !called {
		t.Fatal("Fire did not invoke the tick handler")
	}
}

func TestStep04_Timer_Stop_StopsFiring(t *testing.T) {
	called := false
	tm := scheduler.New(time.Millisecond, func() error {
		called = true
		return nil
	})
	if err := tm.Stop(); err != nil {
		t.Fatalf("Stop: unexpected error: %v", err)
	}
	if !tm.Stopped() {
		t.Fatal("Stopped should be true right after Stop")
	}
	if err := tm.Fire(); err != nil {
		t.Fatalf("Fire after Stop: unexpected error: %v", err)
	}
	if called {
		t.Fatal("Fire after Stop must not invoke the tick handler")
	}
}

func TestStep04_Timer_Interval_IsSettable(t *testing.T) {
	tm := scheduler.New(2500*time.Millisecond, func() error { return nil })
	if got := tm.Interval(); got != 2500*time.Millisecond {
		t.Fatalf("Interval = %v, want %v", got, 2500*time.Millisecond)
	}
}

func TestStep04_Scheduler_Add_Then_RunNext_PopsInFifoOrder(t *testing.T) {
	tm := scheduler.New(time.Millisecond, func() error { return nil })
	tbl := process.New()
	for _, pid := range []int{1, 2, 3} {
		if _, err := tbl.Create(pid, process.Registers{PC: 3}); err != nil {
			t.Fatalf("Create(%d): unexpected error: %v", pid, err)
		}
	}
	sched := scheduler.NewScheduler(tm, tbl)
	for _, pid := range []int{1, 2, 3} {
		if err := sched.Add(pid); err != nil {
			t.Fatalf("Add(%d): unexpected error: %v", pid, err)
		}
	}
	if got := sched.Len(); got != 3 {
		t.Fatalf("Len = %d, want 3", got)
	}

	for _, wantPid := range []int{1, 2, 3} {
		p, err := sched.RunNext()
		if err != nil {
			t.Fatalf("RunNext: unexpected error: %v", err)
		}
		if p.Pid != wantPid {
			t.Fatalf("RunNext returned pid %d, want %d", p.Pid, wantPid)
		}
	}
	if got := sched.Len(); got != 0 {
		t.Fatalf("Len after draining = %d, want 0", got)
	}
	if _, err := sched.RunNext(); !errors.Is(err, scheduler.ErrEmpty) {
		t.Fatalf("RunNext on empty queue error = %v, want an error satisfying errors.Is(err, scheduler.ErrEmpty)", err)
	}
}

func TestStep08_Scheduler_SkipsExitedProcesses(t *testing.T) {
	tm := scheduler.New(time.Millisecond, func() error { return nil })
	tbl := process.New()
	if _, err := tbl.Create(1, process.Registers{PC: 5}); err != nil {
		t.Fatalf("Create(1): unexpected error: %v", err)
	}
	if _, err := tbl.Create(2, process.Registers{PC: 7}); err != nil {
		t.Fatalf("Create(2): unexpected error: %v", err)
	}
	sched := scheduler.NewScheduler(tm, tbl)
	if err := sched.Add(1); err != nil {
		t.Fatalf("Add(1): unexpected error: %v", err)
	}
	if err := sched.Add(2); err != nil {
		t.Fatalf("Add(2): unexpected error: %v", err)
	}
	if _, err := tbl.Exit(1); err != nil {
		t.Fatalf("Exit(1): unexpected error: %v", err)
	}
	p, err := sched.RunNext()
	if err != nil {
		t.Fatalf("RunNext: unexpected error: %v", err)
	}
	if p.Pid != 2 {
		t.Fatalf("RunNext should skip the exited process 1 and run 2, got pid %d", p.Pid)
	}
}

func TestStep04_Timer_Fire_RunsTheNextScheduledProcess(t *testing.T) {
	tm := scheduler.New(time.Millisecond, func() error { return nil })
	tbl := process.New()
	if _, err := tbl.Create(1, process.Registers{PC: 7}); err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	sched := scheduler.NewScheduler(tm, tbl)
	if err := sched.Add(1); err != nil {
		t.Fatalf("Add: unexpected error: %v", err)
	}
	if err := tm.Fire(); err != nil {
		t.Fatalf("Fire: unexpected error: %v", err)
	}
	if got := sched.Len(); got != 0 {
		t.Fatalf("Len after one Fire = %d, want 0 (the single queued process was run)", got)
	}
}
