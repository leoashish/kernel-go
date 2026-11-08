package syscall_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"kernel/internal/syscall"
)

// Step 02 -- the syscall dispatch table. Run: go test ./internal/syscall -run Step02 -v

func TestStep02_Dispatch_CallsTheRegisteredHandler_AndReturnsItsError(t *testing.T) {
	table := syscall.New()
	var gotArgs []any
	h := func(args ...any) error {
		gotArgs = args
		return nil
	}
	table.Register(777, h)

	if err := table.Dispatch(777, 10, 20); err != nil {
		t.Fatalf("Dispatch: unexpected error: %v", err)
	}
	if len(gotArgs) != 2 {
		t.Fatalf("handler got %d args, want 2 (got %v)", len(gotArgs), gotArgs)
	}
	if got := gotArgs[0].(int); got != 10 {
		t.Fatalf("handler arg 0 = %v, want 10", gotArgs[0])
	}
	if gotArgs[1] != 20 {
		t.Fatalf("handler arg 1 = %v, want 20", gotArgs[1])
	}
}

func TestStep02_Dispatch_ReturnsErrNoSys_ForAnUnknownSyscall(t *testing.T) {
	table := syscall.New()

	err := table.Dispatch(4242, 1, 2)
	if !errors.Is(err, syscall.ErrNoSys) {
		t.Fatalf("Dispatch(4242) error = %v, want an error satisfying errors.Is(err, syscall.ErrNoSys)", err)
	}
}

func TestStep02_Lookup_FindsAnExistingHandler_AndMissesForAnUnknownNumber(t *testing.T) {
	table := syscall.New()
	if _, ok := table.Lookup(1); ok {
		t.Fatal("Lookup of an unregistered number should miss")
	}

	table.Register(1, func(args ...any) error { return nil })
	if _, ok := table.Lookup(1); !ok {
		t.Fatal("Lookup of a registered number should hit")
	}
	if _, ok := table.Lookup(2); ok {
		t.Fatal("Lookup of an unregistered number should miss")
	}
}

func TestStep02_Write_WritesTheLineToTheConsole(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	table := syscall.New()
	table.Register(syscall.SYS_WRITE, syscall.Write)

	if err := table.Dispatch(syscall.SYS_WRITE, syscall.StdoutFd, "hello from the kernel"); err != nil {
		t.Fatalf("Dispatch(SYS_WRITE): unexpected error: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	if got := string(data); got != "hello from the kernel\n" {
		t.Fatalf("console output = %q, want %q", got, "hello from the kernel\n")
	}
}

func TestStep02_Write_ReturnsErrBadFd_ForAnUnknownFd(t *testing.T) {
	table := syscall.New()
	table.Register(syscall.SYS_WRITE, syscall.Write)

	err := table.Dispatch(syscall.SYS_WRITE, 99, "hello")
	if !errors.Is(err, syscall.ErrBadFd) {
		t.Fatalf("Dispatch(SYS_WRITE, fd=99) error = %v, want an error satisfying errors.Is(err, syscall.ErrBadFd)", err)
	}
}

func TestStep02_Exit_ReturnsAnExitErrorWithTheCode(t *testing.T) {
	table := syscall.New()
	table.Register(syscall.SYS_EXIT, syscall.Exit)

	err := table.Dispatch(syscall.SYS_EXIT, 7)
	var exitErr *syscall.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Dispatch(SYS_EXIT) error = %v, want an error of type *syscall.ExitError", err)
	}
	if exitErr.Code != 7 {
		t.Fatalf("ExitError.Code = %d, want 7", exitErr.Code)
	}
}
