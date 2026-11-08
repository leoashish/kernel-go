# Step 04 — The timer interrupt and the scheduler's run queue

| | |
|---|---|
| **Build** | `internal/scheduler/scheduler.go` — `Timer`, `Scheduler` |
| **Tests** | `internal/scheduler/scheduler_test.go` — 5 tests |
| **Run** | `go test ./internal/scheduler -run Step04 -v` |
| **You'll learn** | timer interrupts, tick handlers, run queues, FIFO scheduling, `errors.Is` |
| **Time** | ~45–60 min |

## 1. The story

A CPU without a timer just sits there waiting for interrupts that may never come; without a scheduler there is no
"what runs next". Step 04 wires the two: a **timer** that fires interrupts periodically, and a **scheduler** whose tick
handler runs the next process out of a FIFO run queue.

## 2. Concept corner

- **An interrupt is just a function call** in this miniature: the timer's `Fire` invokes the tick handler (the
  scheduler's `RunNext`).
- **The run queue is FIFO**: `Add(pid)` appends, the tick pops from the front.
- **Empty queue**: a tick with nothing queued is not a crash — the kernel treats it as a spurious interrupt.

## 3. Requirements

| Piece | Behaviour |
|---|---|
| `Timer.New(interval, tick)` | timer with period + tick handler. |
| `Timer.Fire()` | invoke tick unless stopped; nil error when stopped. |
| `Timer.Stop()` / `Timer.Stopped()` | stop the timer; later Fires become no-ops. |
| `Timer.Interval()` | the period. |
| `Scheduler.NewScheduler(timer, procs)` | scheduler bound to a timer + process table; its tick = run next. |
| `Scheduler.Add(pid)` | enqueue. |
| `Scheduler.Len()` | queued count. |
| `Scheduler.RunNext()` | pop front → return process; `ErrEmpty` when nothing queued. |

## 4. What the tests pin down

- `Fire` invokes the tick handler.
- `Stop` stops future Fires.
- Interval is settable.
- Add → RunNext pops in FIFO order; drain → `ErrEmpty` on the next pop.
- `Timer.Fire` (via the scheduler tick) runs the next scheduled process.

## 5. Hints

<details><summary>Hint 1 — the timer</summary>

`Fire` checks `stopped` under a mutex, then calls `t.tick()`.

</details>

<details><summary>Hint 2 — the scheduler</summary>

`t.tick = func() error { _, err := s.RunNext(); return err }` in the constructor.

</details>

<details><summary>Hint 3 — FIFO pop</summary>

Pop `queue[0]` and slice the rest; empty → `ErrEmpty`.

</details>

## 6. Stretch challenges

1. Round-robin: when a process is "done running" it goes back to the queue tail (here: re-`Add` after the run).
2. Real `time.Ticker` behind the timer abstraction.

## 7. Reflect

- Why does the interrupt handler (timer tick) live in the scheduler?
- FIFO vs. priority vs. CFS — where does Linux actually differ?

## 8. Done when

`go test ./internal/scheduler -run Step04 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, compare with `_solutions/step04/…`.

## 9. Next

Step 05 — the context switch: save/restore the CPU context when the scheduler moves from one process to the next.
