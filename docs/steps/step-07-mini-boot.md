# Step 07 — The mini boot: wiring the kernel together

| | |
|---|---|
| **Build** | `internal/kernel/kernel.go` — the `Kernel` that wires console + syscalls + scheduler + processes |
| **Tests** | `internal/kernel/kernel_test.go` — 1 test |
| **Run** | `go test ./internal/kernel -run Step07 -v` |
| **You'll learn** | how the pieces fit: console → syscall → scheduler → process → context |
| **Time** | ~30 min |

## 1. The story

Every piece now exists: console, syscall dispatch, processes, timer/scheduler, context, allocator. The last step
connects them into a **mini boot**: a kernel that registers its built-ins, spawns an initial user process, arms the
timer, and lets the timer interrupt drive execution — with the first user action being a `write` to the console.

## 2. Concept corner

- **The kernel is just composition**: no new kernel logic here, only wiring. The magic is that the composition already
  works because each piece was tested in isolation.
- **Boot is: register built-ins → spawn process → arm timer → run (interrupt → dispatch)**.

## 3. Requirements

| Piece | Behaviour |
|---|---|
| `New()` | build console (stdout/info), syscall table, scheduler+timer, process table. |
| `Boot()` | register `SYS_WRITE`/`SYS_EXIT`, create initial process, add to scheduler, run. |
| `Run()` | fire timer → tick runs next process → dispatch `SYS_WRITE` with the boot line. |

## 4. What the test pins down

- `Boot` writes `BootLine + "\n"` to the console (via `os.Stdout` capture).

## 5. Reflect

- Which piece was the hardest to wire? Why?
- What would a real kernel's boot sequence add that this miniature doesn't (real hardware, SMP, …)?

## 6. Done when

`go test ./internal/kernel -run Step07 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, compare with
`_solutions/step07/…`.

## 9. Next

**The arc is complete**: console → syscall → process → scheduler/timer → context → allocator → mini boot. You now have
a miniature Linux built from first principles, one tested piece at a time.
