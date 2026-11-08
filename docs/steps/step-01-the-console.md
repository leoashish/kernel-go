# Step 01 — The kernel console: `Write`, `Writef`

| | |
|---|---|
| **Build** | `internal/console/console.go` — the console device |
| **Tests** | `internal/console/console_test.go` — 5 tests |
| **Run** | `go test ./internal/console -run Step01 -v` |
| **You'll learn** | log levels, minimum-level filtering, `io.Writer`, error propagation |
| **Time** | ~30–45 min |

## 1. The story

The Linux kernel never talks to a screen directly. It has a **console device**: an output sink that formats lines and
hands them to drivers (UART, VGA, serial, `dmesg` pipe). Before anything else — before syscalls, interrupts, drivers,
schedulers — the kernel needs a way to say "something happened".

Step 01 builds that console. It is deliberately tiny: a minimum **level**, an **output writer**, and two methods —
`Write` (a ready string) and `Writef` (a `fmt`-formatted string, the kernel's `printk`-flavoured cousin). Everything
later prints through it.

## 2. Concept corner

- **A log level is just an int.** Real kernels have `KERN_DEBUG`, `KERN_INFO`, `KERN_WARNING`, `KERN_ERROR`; ours are
  `Debug(0) < Info(1) < Warn(2) < Error(3)`.
- **The minimum level is a threshold.** A console configured with `Warn` prints warnings and errors, but not debug or
  info chatter. In the real world this is the `-debug` / `dmesg` level knob; in ours it is a field on the struct.
- **Filtering happens before touching the sink.** Dropping a line is a no-op that returns `nil` — a dropped message
  must not surface as an error.
- **Errors propagate.** If the sink (`io.Writer`) fails, the error comes back to the caller. The console does not swallow
  or translate them (yet).

## 3. Requirements

Implement the console:

| Method | Behaviour |
|---|---|
| `Write(level Level, msg string) error` | If `level < minLevel`: return `nil` (drop). Otherwise write `msg + "\n"` to the sink and return the sink's error. |
| `Writef(level Level, format string, args ...any) error` | Same, but first build the line with `fmt.Sprintf(format, args...)`. |

- `New(out io.Writer, minLevel Level) *Console` returns a working console.
- Lines at or above the minimum are written; lines below are dropped silently.

## 4. What the tests pin down

- `Writef` formats and writes the line (`"boot %s"` + `"kernel"` → `"boot kernel\n"`).
- `Write` writes the message as a one-line (`"hello\n"`).
- Lines **below** the minimum are dropped (debug/info hidden when the minimum is `Warn`).
- An **error-level** line is written even when the minimum is higher than the error (errors always surface).
- The **sink error** propagates out of `Writef`/`Write`.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — the shape of a write</summary>

Compare `level < c.minLevel` first; if true, return `nil`. Then format-and-emit: `fmt.Fprintln(c.out, msg)` or
`fmt.Fprintf(c.out, format+"\n", args...)`, and return the error.

</details>

<details><summary>Hint 2 — where the newline comes from</summary>

Decide whether the newline is added by the console or expected in the format. Pick one and make the tests agree.

</details>

<details><summary>Hint 3 — comparing levels</summary>

`level` and `c.minLevel` are `Level` ints; `level >= c.minLevel` is the emit condition.

</details>

## 6. Stretch challenges

1. Add kernel-flavoured convenience wrappers: `Debugf`, `Infof`, `Warnf`, `Errorf` — each delegating to `Writef`.
2. Add a `Printk(format string, args ...any) error` alias for `Writef(Error, …)` — the real kernel's printk.
3. Make the console safe for concurrent use (interrupt context + process context writing at the same time).

## 7. Reflect

- Why is a "minimum level" a threshold and not a per-line decision made by the caller?
- Why must a dropped line return `nil` and not an error?
- Why wrap formatting (`Writef`) instead of making callers pre-format strings? (Kernel: printk has a format string.)

## 8. Done when

`go test ./internal/console -run Step01 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step01/…`.

## 9. Next

Step 02 — the **syscall dispatch table**: the `write` syscall becomes a call into the console (`Writef(Info, "syscall write ...")`),
and an unknown syscall number becomes `ENOSYS`. Its PRD will be written when you finish this one.
