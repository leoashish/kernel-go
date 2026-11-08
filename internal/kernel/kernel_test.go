package kernel_test

import (
	"io"
	"os"
	"testing"

	"kernel/internal/kernel"
)

// Step 07 -- the mini boot: wire console + syscall + scheduler + processes into one kernel. Run: go test ./internal/kernel -run Step07 -v

func TestStep07_Kernel_Boot_WritesTheBootLineToTheConsole(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	k := kernel.New()
	if err := k.Boot(); err != nil {
		t.Fatalf("Boot: unexpected error: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	if got := string(data); got != kernel.BootLine+"\n" {
		t.Fatalf("console output = %q, want %q", got, kernel.BootLine+"\n")
	}
}
