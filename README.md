# Go Concurrency: 3-Stage Pipeline

This project implements a **3-stage concurrent pipeline** in Go, using **goroutines** and **channels** to coordinate multiple producers and consumers without shared memory.

---

##  Pipeline Design

###  Stage 1: Two Producers
- `Producer 1`: Generates odd numbers from 1 to 29.
- `Producer 2`: Generates even numbers from 2 to 30.
- Both send values to a shared buffered channel `inCh` (capacity: 5).
- Each producer sleeps for a random time between 0–1500ms between sends.

###  Stage 2: Two Consumers
- Each consumer reads from `inCh`, squares the number, and sends it to `outCh` (buffered, capacity: 5).
- Each simulates computation delay by sleeping between 0–3s.

###  Stage 3: Final Filter
- A single goroutine reads from `outCh`.
- It prints the **first** value.
- It only prints a value if it’s **strictly greater than the last printed**.
- Outputs a **strictly increasing** subsequence of perfect squares.

---

##  How to Run

```bash
go run pipeline.go
```

###  Sample Output

(Note: Output varies due to random sleep timing)

```
1
4
9
16
36
49
64
81
```

---

##  Questions for Discussion

### (1) Why will there be no deadlock?

Deadlock is avoided because:

- All channels are **buffered** (capacity = 5), which prevents immediate blocking on send.
- Producers and consumers run in separate goroutines and exit cleanly after completion.
- The `WaitGroup`s ensure that stages signal downstream processes when done:
  - After both producers finish, `inCh` is closed.
  - After consumers finish reading from `inCh`, they close `outCh`.
  - The final stage finishes when `outCh` is closed and drained.
- All goroutines terminate properly, and no one blocks waiting for a message that will never arrive.

### (2) How would this work in Elixir?

In Elixir:

- Each process has a **private mailbox**, so a channel-like construct (multiple writers, multiple readers) must be **simulated**.
- To implement this fan-out/fan-in model:
  - You’d spawn **two producer processes** that send to a **dispatcher process** or use a `Registry`-based router to distribute messages to multiple consumers.
  - The dispatcher (or direct sends) would need to alternate or randomly choose between consumer PIDs.
  - For consumers to send results to a shared stage-3 process, they can send messages to a common **PID** (e.g. `final_filter_pid`), effectively using it as a “channel.”

This design is more manual in Elixir than Go, since mailbox semantics are 1:1 and channels are not native first-class citizens.

---

##  Author

**Nandan Mogili**


