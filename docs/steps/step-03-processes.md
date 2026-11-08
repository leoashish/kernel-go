# Step 03 — Processes: process table, `fork`, `exec`

| | |
|---|---|
| **Build** | `internal/process/process.go` — the process table + fork/exec |
| **Tests** | `internal/process/process_test.go` — 7 tests |
| **Run** | `go test ./internal/process -run Step03 -v` |
| **You'll learn** | process tables, fork semantics (copy registers, new pid), exec (new address space → new PC) |
| **Time** | ~60–90 min |

## 1. The story

The kernel is useless without tasks to run. Step 03 gives the kernel a **process table** and the two syscalls that
create/destroy execution contexts: `fork` (copy an existing process) and `exec` (replace a process's address space).

## 2. Concept corner

- **A process = registers + state** in this miniature: pid, state, PC + GPRs. Real kernels add mm, file tables, etc.;
  here the copy semantics are what matters.
- **`fork`**: child gets a **new pid** and **copies the parent's registers**; child starts `Ready`.
- **`exec`**: replaces the address space → the program counter moves to the new entry point; process becomes `Ready`.

## 3. Requirements

| Piece | Behaviour |
|---|---|
| `Create(pid, regs)` | add process; error `ErrProcessExists` if pid taken. |
| `Fork(pid, parent)` | child with `pid` copying `parent`'s regs, state `Ready`; errors `ErrNoSuchProcess` (nil/unknown parent), `ErrProcessExists`. |
| `Get(pid)` | return process or error `ErrNoSuchProcess`. |
| `Exec(pid, entry)` | set regs.PC = entry, state `Ready`; error `ErrNoSuchProcess` if unknown. |

All mutations guarded by the table mutex; safe for concurrent use.

## 4. What the tests pin down

- `Create` stores regs and reports the state `Ready`.
- `Create` twice with the same pid → `ErrProcessExists`.
- `Fork` copies parent regs into a child with a different pid, state `Ready`.
- `Fork` with nil parent → `ErrNoSuchProcess`.
- `Get` of unknown pid → `ErrNoSuchProcess`.
- `Exec` moves the PC to the entry and resets state to `Ready`.
- `Exec` of unknown pid → `ErrNoSuchProcess`.

## 5. Hints

<details><summary>Hint 1 — the table</summary>

Map[pid]*Process under a sync.Mutex (same shape as the console's sink decision: one writer).

</details>

<details><summary>Hint 2 — fork</summary>

Copy `parent.Regs` (value copy) into the child; child `State = Ready`.

</details>

<details><summary>Hint 3 — exec</summary>

Set `p.Regs.PC = entry` and `p.State = Ready`.

</details>

## 6. Stretch challenges

1. Add `Kill(pid)` — mark `Exited` (the scheduler stops queueing exited processes).
2. Add `List() []int` — live pids, for debugging/procfs-style introspection.

## 7. Reflect

- Why does fork copy the whole address space vs. exec replacing it?
- Why does the child start `Ready` (not `Running`)?

## 8. Done when

`go test ./internal/process -run Step03 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, compare with `_solutions/step03/…`.

## 9. Next

Step 04 — the scheduler + timer interrupt: the run queue that ties processes to execution.
