// REFERENCE SOLUTION for Step 06. Try it yourself first!
// Package memory implements the kernel's heap allocator: a free-list allocator over a fixed-size heap.
package memory

import (
	"errors"
	"sort"
	"sync"
)

// ErrOutOfMemory: the heap does not have a free block big enough for the allocation. GIVEN.
var ErrOutOfMemory = errors.New("memory: out of memory")

// Block is one free (or allocated) region of the heap. GIVEN.
type Block struct {
	Addr int
	Size int
}

// Allocator hands out memory from a fixed-size heap using a free list. GIVEN.
type Allocator struct {
	mu   sync.Mutex
	heap int
	free []Block
}

// New returns an allocator over a heap of the given size (addressed from 0). GIVEN.
func New(heapSize int) *Allocator {
	return &Allocator{heap: heapSize, free: []Block{{Addr: 0, Size: heapSize}}}
}

// Alloc takes the lowest free block that fits size bytes and returns its address.
func (a *Allocator) Alloc(size int) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := range a.free {
		b := &a.free[i]
		if b.Size >= size {
			addr := b.Addr
			b.Addr += size
			b.Size -= size
			if b.Size == 0 {
				a.free = append(a.free[:i], a.free[i+1:]...)
			}
			return addr, nil
		}
	}
	return 0, ErrOutOfMemory
}

// Free returns size bytes at addr to the free list (coalescing with adjacent blocks).
func (a *Allocator) Free(addr int, size int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.free = append(a.free, Block{Addr: addr, Size: size})
	sort.Slice(a.free, func(i, j int) bool { return a.free[i].Addr < a.free[j].Addr })
	merged := make([]Block, 0, len(a.free))
	for _, b := range a.free {
		if len(merged) > 0 && merged[len(merged)-1].Addr+merged[len(merged)-1].Size == b.Addr {
			merged[len(merged)-1].Size += b.Size
			continue
		}
		merged = append(merged, b)
	}
	a.free = merged
	return nil
}
