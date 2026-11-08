# Step 02 — The syscall dispatch table: `write`, `exit`, `ENOSYS`

| | |
|---|---|
| **Build** | `internal/syscall/syscall.go` — the dispatch table + built-in handlers |
| **Tests** | `internal/syscall/syscall_test.go` — 6 tests |
| **Run** | `go test ./internal/syscall -run Step02 -v` |
| **You'll learn** | dispatch tables, `ENOSYS`/`EBADF` errnos, handler wiring, `errors.Is`/`errors.As` |
| **Time** | ~45–60 min |

## 1. The story

The kernel and user programs speak different languages. User programs don't execute kernel code directly — they make
**syscalls**: numbered requests (`write`, `exit`, `read`, …) that enter the kernel through a fixed door: the dispatch
table. The kernel owns the table; user code gets a number and arguments; the table routes to a handler.

Step 02 builds that table. `Dispatch(syscallNumber, args...)` is the single entry point the future "process" layer will
call; the built-in handlers do the real work: `write` writes to the **console** (Step 01), `exit` reports termination.

## 2. Concept corner

- **A dispatch table is a map from number → function.** Simple, idiomatic Go; real kernels use tables or IDT-like
  structures (x86 Interrupt Descriptor Table is the hardware cousin of this).
- **`ENOSYS`** is the error a kernel returns for a syscall it does not implement yet. It is not a panic — it is a
  graceful "not supported, here is your error code".
- **`EBADF`** is the error for a bad file descriptor. Our miniature has fds `1` (stdout) and `2` (stderr), both mapped to
  the console.
- **Errors are sentinels or typed values**: `ErrNoSys` is a sentinel (`errors.Is`); `ExitError` is a typed error
  (`errors.As`) that carries the exit code.

## 3. Requirements

Implement in `internal/syscall`:

| Piece | Behaviour |
|---|---|
| `New() *DispatchTable` | empty table |
| `Register(num int, h Handler)` | install/replace the handler for `num` |
| `Lookup(num int) (Handler, bool)` | find a handler; report presence |
| `Dispatch(num int, args ...any) error` | look up + invoke; unknown number → `ErrNoSys` |
| `Write(args ...any) error` | `SYS_WRITE` handler: args `(fd, line)`; fd must be `StdoutFd`/`StderrFd` else `ErrBadFd`; write line to console; nil on success |
| `Exit(args ...any) error` | `SYS_EXIT` handler: return an `*ExitError` carrying the code |

## 4. What the tests pin down

- Dispatching to a registered handler invokes it with its args and returns its error.
- Dispatching an unknown syscall number returns an error satisfying `errors.Is(err, syscall.ErrNoSys)`.
- `Lookup` hits for registered numbers, misses for unknown ones.
- `Dispatch(SYS_WRITE, StdoutFd, "line")` writes `"line\n"` to the console (via `os.Stdout` capture).
- `Dispatch(SYS_WRITE, 99, "line")` returns an error satisfying `errors.Is(err, syscall.ErrBadFd)`.
- `Dispatch(SYS_EXIT, code)` returns an `*syscall.ExitError` with the right code.

## 5. Hints

<details><summary>Hint 1 — the table</summary>

`New` returns a struct with an unexported `map[int]Handler`; `Register` stores, `Lookup` reads.

</details>

<details><summary>Hint 2 — dispatch</summary>

```go
h, ok := t.Lookup(num)
if !ok { return ErrNoSys }
return h(args...)
```

</details>

<details><summary>Hint 3 — write to the console</summary>

`fd` must be `StdoutFd` (1) or `StderrFd` (2); otherwise return an error satisfying `errors.Is(err, ErrBadFd)`. Then
`console.New(os.Stdout, console.Info).Writef(console.Info, "%s", line)` and return its error.

</details>

<details><summary>Hint 4 — exit</summary>

Return a new `*ExitError{Code: code}` — or wrap it. The process teardown itself is Step 03's job.

</details>

## 6. Stretch challenges

1. Add `SYS_READ`/`SYS_OPEN`-style stubs that return `ErrNoSys` (i.e. register handlers that immediately return `ErrNoSys`).
2. Add `RegisterAll`-style default table wiring the built-ins: `syscall.New()` returns a table with `SYS_WRITE`/`SYS_EXIT` pre-registered.
3. Make the console sink injectable into the write handler (constructor injection) so tests do not swap `os.Stdout`.

## 7. Reflect

- Why is `ENOSYS` the right behaviour for an unimplemented syscall (vs. panicking)?
- Why do real kernels number syscalls and keep a table (vs. switch statements)?
- What would change if syscalls took typed args instead of ints?

## 8. Done when

`go test ./internal/syscall -run Step02 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with
`_solutions/step02/…`.

## 9. Next

Step 03 — processes: the process table that the dispatch table now lives in, plus `fork`/`exec`/`exit` that wire real
syscalls to real processes. Its PRD will be written when you finish this one.
