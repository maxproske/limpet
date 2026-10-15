# P1 · tinybox: Unix tools in Go

Runs on the macOS host. Status is tracked in [`PROGRESS.md`](../../PROGRESS.md).

## What you're building

`tinybox` is one Go binary containing several small Unix tools ("applets"), in the spirit of BusyBox:

```
tinybox echo hello      tinybox cat notes.txt      tinybox env      tinybox probe --json
```

You'll build it one applet at a time. Each applet is small, and each was chosen because it makes you confront a piece of the process model that limpet depends on.

## Why this project comes first

- You learn Go fundamentals on problems that matter later: argv, exit status, the standard streams, the environment, process identity.
- `probe` becomes the course's test subject. From P2 on, we run it inside every new kind of sandbox to see what changed.
- A statically linked Linux build of tinybox can be the *only* program in a root filesystem. You'll see that first in a `FROM scratch` container (L5) and use it for real in P6.

## Lesson plan

Only lessons marked *ready* exist. The rest is a sketch, to be written or changed based on how the earlier lessons go.

| # | Lesson | Status | Core ideas |
|---|---|---|---|
| 1 | [Hello, process](lessons/01-hello-process.md) | ready | toolchain and modules; argv; exit status; `true`, `false`, `echo` |
| 2 | [Streams and file descriptors](lessons/02-streams-and-fds.md) | ready | fds 0/1/2; redirection and pipes; EOF and SIGPIPE; `io.Reader`/`io.Writer`; errors; `cat` |
| 3 | Testable tools | sketch | one applet signature with injected stdin/stdout/stderr; table-driven tests; `flag.FlagSet` |
| 4 | Who am I? | sketch | the environment; pid, ppid, uid, gid, cwd, hostname; structs and JSON; `env` and `probe` |
| 5 | One binary, many names | sketch | `argv[0]` dispatch through symlinks; static vs. dynamic linking; cross-compiling; tinybox alone in a `FROM scratch` container |

## Done when

- Every applet behaves correctly with pipes, redirection, bad input, and exit statuses.
- Each applet is tested through injected stdin/stdout/stderr (from L3 on).
- Without notes, you can explain who builds argv, what an exit status can and can't carry, and what fds 0/1/2 are and how the shell wires them up.

## Extensions

`wc`, `head`, `tee`, `yes` (make it fast), `cat -n`, `seq`.

## Where the code goes

This folder becomes a Go module when you run `go mod init` in Lesson 1. Build output goes to `bin/`, which git ignores:

```
cd projects/01-tinybox
go build -o bin/ .     # writes bin/tinybox
go run . echo hi       # builds to a temporary directory and runs it in one step
```
