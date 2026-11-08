package memory_test

import (
	"errors"
	"testing"

	"kernel/internal/memory"
)

// Step 06 -- the kernel heap allocator. Run: go test ./internal/memory -run Step06 -v

func TestStep06_Alloc_ReturnsDistinctAddresses_InIncreasingOrder(t *testing.T) {
	a := memory.New(100)
	first, err := a.Alloc(10)
	if err != nil {
		t.Fatalf("Alloc(10): unexpected error: %v", err)
	}
	second, err := a.Alloc(10)
	if err != nil {
		t.Fatalf("Alloc(10): unexpected error: %v", err)
	}
	if first != 0 || second != 10 {
		t.Fatalf("Alloc returned (%d, %d), want (0, 10)", first, second)
	}
}

func TestStep06_Alloc_ReturnsErrOutOfMemory_WhenTheHeapIsTooSmall(t *testing.T) {
	a := memory.New(10)
	if _, err := a.Alloc(20); !errors.Is(err, memory.ErrOutOfMemory) {
		t.Fatalf("Alloc(20) on a 10-byte heap error = %v, want an error satisfying errors.Is(err, memory.ErrOutOfMemory)", err)
	}
}

func TestStep06_Free_MakesMemoryAvailableAgain(t *testing.T) {
	a := memory.New(100)
	addr, err := a.Alloc(10)
	if err != nil {
		t.Fatalf("Alloc(10): unexpected error: %v", err)
	}
	if addr != 0 {
		t.Fatalf("first Alloc = %d, want 0", addr)
	}
	if err := a.Free(0, 10); err != nil {
		t.Fatalf("Free(0, 10): unexpected error: %v", err)
	}
	addr2, err := a.Alloc(10)
	if err != nil {
		t.Fatalf("Alloc(10) after Free: unexpected error: %v", err)
	}
	if addr2 != 0 {
		t.Fatalf("Alloc after Free = %d, want 0 (the memory came back)", addr2)
	}
}

func TestStep06_Free_CoalescesAdjacentBlocks(t *testing.T) {
	a := memory.New(100)
	addr1, err := a.Alloc(10) // 0
	if err != nil {
		t.Fatalf("Alloc(10): unexpected error: %v", err)
	}
	addr2, err := a.Alloc(10) // 10
	if err != nil {
		t.Fatalf("Alloc(10): unexpected error: %v", err)
	}
	if addr1 != 0 || addr2 != 10 {
		t.Fatalf("Alloc returned (%d, %d), want (0, 10)", addr1, addr2)
	}
	if err := a.Free(addr1, 10); err != nil {
		t.Fatalf("Free: unexpected error: %v", err)
	}
	if err := a.Free(addr2, 10); err != nil {
		t.Fatalf("Free: unexpected error: %v", err)
	}
	addr3, err := a.Alloc(20) // coalesced block starting at 0
	if err != nil {
		t.Fatalf("Alloc(20) after freeing both halves: unexpected error: %v", err)
	}
	if addr3 != 0 {
		t.Fatalf("Alloc(20) after coalescing = %d, want 0", addr3)
	}
}
