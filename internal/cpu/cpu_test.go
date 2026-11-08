package cpu_test

import (
	"testing"

	"kernel/internal/cpu"
	"kernel/internal/process"
)

// Step 05 -- context switches: save/restore of the CPU context. Run: go test ./internal/cpu -run Step05 -v

func TestStep05_Save_ReturnsTheCurrentContext(t *testing.T) {
	proc := &process.Process{Pid: 1, State: process.Ready, Regs: process.Registers{PC: 5}}
	cpu.Restore(cpu.Context{Process: proc, Regs: process.Registers{PC: 5}})

	got := cpu.Save()
	if got.Process != proc {
		t.Fatalf("Save Process = %v, want %v", got.Process, proc)
	}
	if got.Regs.PC != 5 {
		t.Fatalf("Save Regs.PC = %d, want 5", got.Regs.PC)
	}
}

func TestStep05_Restore_MakesTheContextCurrent(t *testing.T) {
	proc := &process.Process{Pid: 2, State: process.Ready, Regs: process.Registers{PC: 9}}
	ctx := cpu.Context{Process: proc, Regs: process.Registers{PC: 9}}
	cpu.Restore(ctx)

	if got := cpu.Save(); got.Process != proc || got.Regs.PC != 9 {
		t.Fatalf("after Restore, Save = %+v, want the context %+v", got, ctx)
	}
}

func TestStep05_Switch_SavesThePreviousContextAndLoadsTheNewOne(t *testing.T) {
	a := cpu.Context{Process: &process.Process{Pid: 1, State: process.Ready, Regs: process.Registers{PC: 5}}, Regs: process.Registers{PC: 5}}
	b := cpu.Context{Process: &process.Process{Pid: 2, State: process.Ready, Regs: process.Registers{PC: 7}}, Regs: process.Registers{PC: 7}}
	cpu.Restore(a)

	saved := cpu.Switch(b)
	if saved.Process != a.Process || saved.Regs.PC != a.Regs.PC {
		t.Fatalf("Switch saved the wrong context: %+v, want %+v", saved, a)
	}
	if got := cpu.Save(); got.Process != b.Process || got.Regs.PC != b.Regs.PC {
		t.Fatalf("after Switch, current context = %+v, want %+v", got, b)
	}
}

func TestStep05_Switch_Then_Restore_GoesBack(t *testing.T) {
	a := cpu.Context{Process: &process.Process{Pid: 1, State: process.Ready, Regs: process.Registers{PC: 5}}, Regs: process.Registers{PC: 5}}
	b := cpu.Context{Process: &process.Process{Pid: 2, State: process.Ready, Regs: process.Registers{PC: 7}}, Regs: process.Registers{PC: 7}}
	cpu.Restore(a)
	saved := cpu.Switch(b)
	cpu.Restore(saved)

	if got := cpu.Save(); got.Process != a.Process || got.Regs.PC != a.Regs.PC {
		t.Fatalf("after Switch+Restore, current context = %+v, want %+v", got, a)
	}
}
