package process_test

import (
	"errors"
	"testing"

	"kernel/internal/process"
)

// Step 03 -- processes: table, fork, exec. Run: go test ./internal/process -run Step03 -v

func TestStep03_Create_ReturnsTheProcessWithTheGivenRegs(t *testing.T) {
	tbl := process.New()

	p, err := tbl.Create(1, process.Registers{PC: 5})
	if err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	if p.Pid != 1 {
		t.Fatalf("Create pid = %d, want 1", p.Pid)
	}
	if p.Regs.PC != 5 {
		t.Fatalf("Create regs.PC = %d, want 5", p.Regs.PC)
	}
	if p.State != process.Ready {
		t.Fatalf("Create state = %v, want %v", p.State, process.Ready)
	}
}

func TestStep03_Create_ReturnsErrForAnExistingPid(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Create(1, process.Registers{PC: 1}); err != nil {
		t.Fatalf("first Create: unexpected error: %v", err)
	}
	if _, err := tbl.Create(1, process.Registers{PC: 2}); !errors.Is(err, process.ErrProcessExists) {
		t.Fatalf("second Create error = %v, want an error satisfying errors.Is(err, process.ErrProcessExists)", err)
	}
}

func TestStep03_Fork_CopiesTheParentsRegsIntoANewChild(t *testing.T) {
	tbl := process.New()
	parent, err := tbl.Create(1, process.Registers{PC: 5})
	if err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	child, err := tbl.Fork(2, parent)
	if err != nil {
		t.Fatalf("Fork: unexpected error: %v", err)
	}
	if child.Pid == parent.Pid {
		t.Fatalf("Fork child pid = %d, want a different pid than the parent %d", child.Pid, parent.Pid)
	}
	if child.Regs.PC != parent.Regs.PC {
		t.Fatalf("Fork child regs.PC = %d, want the parent's %d", child.Regs.PC, parent.Regs.PC)
	}
	if child.State != process.Ready {
		t.Fatalf("Fork child state = %v, want %v", child.State, process.Ready)
	}
}

func TestStep03_Fork_ReturnsErrNoSuchProcess_ForUnknownParent(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Fork(2, nil); !errors.Is(err, process.ErrNoSuchProcess) {
		t.Fatalf("Fork with nil parent error = %v, want an error satisfying errors.Is(err, process.ErrNoSuchProcess)", err)
	}
}

func TestStep03_Get_ReturnsErrNoSuchProcess_ForUnknownPid(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Get(99); !errors.Is(err, process.ErrNoSuchProcess) {
		t.Fatalf("Get(99) error = %v, want an error satisfying errors.Is(err, process.ErrNoSuchProcess)", err)
	}
}

func TestStep03_Exec_SetsTheProgramsCounter(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Create(1, process.Registers{PC: 5}); err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	if _, err := tbl.Exec(1, 42); err != nil {
		t.Fatalf("Exec: unexpected error: %v", err)
	}
	got, err := tbl.Get(1)
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if got.Regs.PC != 42 {
		t.Fatalf("after Exec regs.PC = %d, want 42", got.Regs.PC)
	}
	if got.State != process.Ready {
		t.Fatalf("after Exec state = %v, want %v", got.State, process.Ready)
	}
}

func TestStep03_Exec_ReturnsErrNoSuchProcess_ForUnknownPid(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Exec(99, 42); !errors.Is(err, process.ErrNoSuchProcess) {
		t.Fatalf("Exec(99) error = %v, want an error satisfying errors.Is(err, process.ErrNoSuchProcess)", err)
	}
}

func TestStep08_Process_Exit_MarksTheProcessExited(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Create(1, process.Registers{PC: 5}); err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	if _, err := tbl.Exit(1); err != nil {
		t.Fatalf("Exit: unexpected error: %v", err)
	}
	got, err := tbl.Get(1)
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if got.State != process.Exited {
		t.Fatalf("after Exit state = %v, want %v", got.State, process.Exited)
	}
}

func TestStep08_Process_Exit_ErrorsForUnknownPid(t *testing.T) {
	tbl := process.New()
	if _, err := tbl.Exit(99); !errors.Is(err, process.ErrNoSuchProcess) {
		t.Fatalf("Exit(99) error = %v, want an error satisfying errors.Is(err, process.ErrNoSuchProcess)", err)
	}
}
