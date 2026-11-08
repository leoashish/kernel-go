# Step 08 — Exit and teardown: when a process leaves

| | |
|---|---|
| **Build** | `internal/process` (`Exit`), `internal/scheduler` (skip exited), `internal/kernel` (teardown on exit) |
| **Tests** | `internal/process` +2, `internal/scheduler` +1 — 3 new tests |
| **Run** | `go test ./internal/process -run Step08 -v` · `go test ./internal/scheduler -run Step08 -v` · `go test ./internal/kernel -run Step07 -v` |
| **You'll learn** | process exit states, teardown, scheduler invariants |
| **Time** | ~30–45 min |

## 1. The story

A process that finishes (or crashes) must leave the system cleanly: its state becomes `Exited`, and the scheduler must
never run it again.

## 2. Concept corner

- **`Exit(pid)`** marks a process `Exited` (the process table owns process state).
- **The scheduler skips exited processes** when running the queue (the tick handler is the run-next).
- **The kernel tears down** when the action of a running process ends in an exit.

## 3. Requirements

| Piece | Behaviour |
|---|---|
| `ProcessTable.Exit(pid)` | mark process `Exited`; error `ErrNoSuchProcess` for unknown pid. |
| `Scheduler.RunNext` | skip processes in `Exited` state when popping. |
| `Kernel.executeBoot` | when the executed action returns an `*syscall.ExitError`, mark the running process exited (teardown). |

## 4. What the tests pin down

- `Exit` marks a process `Exited`; `Exit` of unknown pid errors with `ErrNoSuchProcess`.
- `RunNext` skips exited processes and runs the next live one.
- (Kernel) the boot action that ends in an exit triggers teardown — via the marked-exited path.

## 5. Hints

<details><summary>Hint 1 — Exit</summary>

Set `p.State = Exited` on the found process; unknown pid → `ErrNoSuchProcess`.

</details>

<details><summary>Hint 2 — RunNext skip</summary>

Pop, look up, `if p.State == process.Exited { continue }`.

</details>

## 6. Stretch challenges

1. `Kill(pid)` — force-exit a process from outside.
2. Queue exited processes for cleanup/reaping (zombie semantics).

## 7. Done when

`go test ./... -race -count=1` → **all green**.

## 9. Next

The miniature kernel is now complete end to end: console → syscall → process → scheduler/timer → context → memory →
mini boot → **exit/teardown**.
