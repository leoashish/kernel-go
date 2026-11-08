# Kernel (Go) — learn Linux kernel internals by building a miniature kernel from first principles

A hands-on course in **operating systems**, built the same way as [distkv-go](../distkv-go) and
[blockchain-go](../blockchain-go): **one step at a time, each with a PRD, failing tests, stubs, collapsible hints,
stretch challenges and a reference solution**.

Today you are building **Kernel** — a miniature Linux-like kernel — out of first principles, in Go. Step 01 gives the
kernel a console to write boot lines to. Step 02 gives it a **syscall dispatch table** (`write`, `exit`, unknown →
`ENOSYS`). Step 03 gives it **processes** (`fork`, `exec`). Step 04 adds the **timer interrupt and the scheduler**. Step
05 adds **context switches**. Step 06 adds a **kernel heap** allocator. Step 07 wires the **mini boot**: one kernel, end to end.

The goal: every day you add one small, testable piece of the kernel, and you *understand* why it is there — because you
wrote it, tested it red-first, and debugged it.

## What you are building, and why this shape

A real Linux kernel is ~30 million lines you will never read. The ~300 lines of core that every kernel shares are the
valuable ones, and they are small enough to build in a weekend of TDD:

1. **Console** — how does anything get printed? (Step 01.)
2. **Syscall dispatch** — how do user programs ask the kernel for things (`write`, `exit`, …)? (Step 02.)
3. **Processes** — `fork`, `exec`, the process table. (Step 03.)
4. **Timer + scheduler** — what runs when, and in what order? (Step 04.)
5. **Context switch** — how does the CPU move from one task to another? (Step 05.)
6. **Memory allocator** — how does the kernel hand out memory? (Step 06.)
7. **Mini boot** — all of the above wired together into one booting kernel. (Step 07.)

By the end you have built the skeleton of everything Linux does — **so you know where the real kernel gets it right and
where it makes trade-offs you haven't felt yet**.

## Start here

```bash
cd kernel-go
go version                      # 1.22 or newer
go test ./...                   # 32 tests, all passing
go test ./internal/console -run Step01 -v
```

You should see **32 passing tests**, all green — `ok` for every package, no failures, no hangs. That is the "green" in
red-green-refactor. 🎯

Work through a step:

1. Read the step's PRD (`docs/steps/step-01-the-console.md` — the current one; siblings are in `docs/steps/`).
2. Open the tests (`internal/<package>/<step>_test.go`) and read them **before** coding. They are the spec.
3. Open the stub (the doc comment on each stub repeats the requirements).
4. Run only that step while you work: `go test ./internal/console -run Step01 -v`.
5. When green, read **Reflect**, try a **Stretch**, run `go vet ./...` and `gofmt -l .`, then compare with
   `_solutions/stepNN/…`.
6. `git commit`, next step.

**Stuck?** Hints in each PRD go from gentle to explicit. As a last resort:

```bash
bash scripts/apply-solution.sh 01            # macOS / Linux / Git-Bash
.\scripts\apply-solution.ps1 01              # Windows PowerShell
```

> Note: steps build on earlier ones. The syscall step's tests need the console (Step 01) to pass.

## Roadmap

| # | Package | Feature | Go ideas | Patterns / principles | Tests |
|---|---------|---------|----------|-----------------------|:------:|
| 01 | `console` | The kernel console: levels and sinks | `io.Writer`, `fmt`, levels | logging, thresholds | 5 ✅ |
| 02 | `syscall` | The syscall dispatch table: `write`, `exit`, `ENOSYS` | function tables, closures | dispatch, errnos | 6 ✅ |
| 03 | `process` | The kernel process table: `fork`, `exec` | structs, mutexes | fork semantics, tables | 7 ✅ |
| 04 | `scheduler` | The timer interrupt and the scheduler's run queue | `time`, channels, locks | interrupts, queues | 5 ✅ |
| 05 | `cpu` | [Context switches: save/restore registers](docs/steps/step-05-context-switch.md) | stack frames | context switch mechanics | 4 ✅ |
| 06 | `memory` | The kernel heap allocator: free-list allocator | maps, free lists | allocation strategies | 4 ✅ |
| 07 | `kernel` | The mini boot: wiring it all together | composition, wiring | integration | 1 ✅ |
| 08 | `process` / `scheduler` / `kernel` | Exit and teardown: mark exited, skip in queue, tear down | states, queues | process lifecycle | 3 ✅ |

**35 tests in total** — all steps solved (01 console → 08 exit/teardown).

## What you'll learn, in order

1. **Step 01 — the console:** log levels, sinks, `io.Writer`, error propagation. ✅
2. **Step 02 — syscalls:** dispatch tables, `write` → console, `ENOSYS` for unknown syscalls. ✅
3. **Step 03 — processes:** fork/exec, process table, states. ✅
4. **Step 04 — timer + scheduler:** interrupts, ticks, run queues, round robin. ✅
5. **Step 05 — context switch:** save/restore, stack switching. ✅
6. **Step 06 — memory:** kernel allocator, free lists, fragmentation. ✅
7. **Step 07 — mini boot:** wiring console + syscalls + scheduler + processes + context + allocator into one kernel. ✅
8. **Step 08 — exit/teardown:** mark exited, skip in queue, tear down. ✅

## Project layout

```
kernel-go/
├── README.md
├── go.mod
├── docs/
│   └── steps/             # ← the PRDs, one per step
├── internal/
│   ├── todo/              # ErrNotImplemented for stubs
│   ├── console/           # step 01: the kernel console
│   ├── syscall/           # step 02: the syscall dispatch table
│   ├── process/           # step 03: the process table, fork, exec
│   ├── scheduler/         # step 04: timer interrupt and run queue
│   ├── cpu/               # step 05: context switch
│   ├── memory/            # step 06: kernel heap allocator
│   └── kernel/            # step 07: the mini boot
├── _solutions/            # reference solutions (ignored by the Go tool)
│   └── step01/ … step07/
└── scripts/               # apply-solution.sh / apply-solution.ps1
```

**Why `_solutions`?** The Go tool ignores directories whose names start with `_`, so `go test ./...` and `go vet ./...`
never try to compile the reference solutions as if they were packages.

## Ground rules that make this effective

- Don't edit the tests to make them pass. If a test looks wrong, that is a great question — read the PRD again, then ask.
- Resist peeking at `_solutions/` until your tests are green.
- Run `go vet ./...` and `gofmt -l .` before every commit.
- Keep the **Reflect** answers in a notes file. Explaining *why* in your own words is what makes it stick.

## How this project was built

- Starting state: all **35 step tests fail cleanly** — no panics, no hangs, no compile errors in the stub packages.
- Solved state (verified with `_solutions/step01`…`_solutions/step08` applied cumulatively): **all 35 pass under `-race`**
  with `go vet` and `gofmt` clean.

## License

MIT — learn stuff.
