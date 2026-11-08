# Step 06 — The kernel heap allocator: free-list allocator

| | |
|---|---|
| **Build** | `internal/memory/memory.go` — `Allocator`, `Block` |
| **Tests** | `internal/memory/memory_test.go` — 4 tests |
| **Run** | `go test ./internal/memory -run Step06 -v` |
| **You'll learn** | free lists, lowest-fit allocation, coalescing, out-of-memory |
| **Time** | ~45 min |

## 1. The story

The kernel hands out memory constantly (sockets buffers, inode structs, page tables). It cannot just bump-allocate
forever; it uses a **heap allocator** with a free list: keep track of free regions, hand out the lowest fitting block,
and return blocks on free (coalescing neighbours).

## 2. Concept corner

- **Bump allocation** (simple, fast, no free) vs **free-list** (reuses freed memory).
- **Lowest fit**: first free block that fits → simple and deterministic.
- **Coalescing**: adjacent free blocks merge so future allocations are not fragmented.

## 3. Requirements

| Piece | Behaviour |
|---|---|
| `New(heapSize)` | allocator over `heapSize` bytes addressed from 0; one free block. |
| `Alloc(size)` | carve the lowest fitting free block; `ErrOutOfMemory` if none fits. |
| `Free(addr, size)` | return the block to the free list, coalescing adjacent blocks. |

## 4. What the tests pin down

- `Alloc` returns ascending, distinct addresses.
- `Alloc` over the heap → `ErrOutOfMemory`.
- `Free` + `Alloc` reuses the same address.
- Freeing adjacent blocks coalesces so a bigger alloc reuses the space.

## 5. Hints

<details><summary>Hint 1 — Alloc</summary>

Scan the free list for the first block with `Size >= size`; carve; drop the block if it empties.

</details>

<details><summary>Hint 2 — Free</summary>

Insert the block, sort by address, merge neighbours that touch.

</details>

## 6. Stretch challenges

1. Buddy allocator (power-of-two sizes) instead of a flat free list.
2. Track allocated (not free) blocks to validate `Free` addresses.

## 7. Reflect

- Why do real kernels use buddy allocators / slab allocators rather than flat free lists?

## 8. Done when

`go test ./internal/memory -run Step06 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, compare with `_solutions/step06/…`.

## 9. Next

Step 07 — the **mini boot**: wire console + syscall + scheduler + processes + context switch + allocator into one
kernel and watch a `write` syscall flow end to end.
