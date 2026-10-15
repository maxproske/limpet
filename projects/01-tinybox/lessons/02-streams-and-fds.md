# Lesson 2 · Streams and file descriptors

> P1 tinybox · runs on macOS · prerequisites: Lesson 1
>
> **Plan: two sessions.** Session A: the challenge and experiments 1–4. Session B: experiments 5–7 and the understanding checks.

## 1. Objective

By the end of this lesson you will:

- have added a `cat` applet that copies files (or standard input) to standard output, reports problems on standard error, keeps going after a bad file, and exits with the right status;
- be able to explain what a file descriptor is, how the shell sets up `<`, `>`, `2>`, and `|` before your program runs, and why your program usually can't tell whether it's writing to a terminal, a file, or a pipe;
- know Go's `io.Reader` contract well enough to write a correct copy loop by hand.

## 2. Why it matters

Containers are mostly plumbing. When limpet starts a process, it decides where that process's input comes from and where its output and errors go: a terminal, a log file, or a pipe back to limpet. Ever wondered why `docker logs` can keep stdout and stderr apart, but stops being able to once you add `-t`? Without a TTY, the runtime gives the container's fd 1 and fd 2 *separate* pipes. With `-t`, both point at one terminal device, so the streams merge before Docker ever sees them.

Files, pipes, terminals, sockets, `/dev/null`, and later the special files under `/proc` and `/sys` are all reached through file descriptors. This lesson builds the model you'll use every time limpet wires up a process.

## 3. The mental model

### File descriptors

For each process, the kernel keeps a **file descriptor table**: small integers that refer to things the process has open. When you open a file, the kernel creates an **open file description** (which object, whether it was opened for reading or writing, and the current position) and hands the process the lowest unused number that points to it. Your program only ever holds the number.

```
process: tinybox cat a.txt > out.txt

fd table (per process)    open file descriptions (kernel)        object
┌───┐
│ 0 │ ─────────────────▶  terminal, read+write              ──▶  /dev/ttys003
│ 1 │ ─────────────────▶  out.txt, write-only, position 0   ──▶  out.txt
│ 2 │ ─────────────────▶  terminal, read+write              ──▶  /dev/ttys003
│ 3 │ ─────────────────▶  a.txt, read-only, position 0      ──▶  a.txt
└───┘
```

The picture simplifies one thing: fds 0 and 2 usually share a *single* open file description, because the terminal was opened once and the descriptor duplicated. We'll need that distinction in P2. Also, `a.txt` gets the lowest free number, usually 3, but the Go runtime sometimes opens descriptors of its own.

### The 0/1/2 convention

Nothing in the kernel makes fd 1 special. It's a convention every Unix program follows: read input from fd 0, write output to fd 1, write diagnostics to fd 2. That convention is what makes programs composable.

### Redirection and pipes happen before your program runs

When you type `tinybox cat a.txt > out.txt`, zsh forks. In the child, *before* exec'ing tinybox, it opens `out.txt` (creating or truncating it) and installs it as fd 1. Then it execs. tinybox inherits that fd table, starts with fd 1 already pointing at the file, and writes there without knowing.

A **pipe** is a kernel buffer with a write end and a read end. For `tinybox cat a.txt | wc -l`, the shell creates the pipe, starts two children, and makes the write end fd 1 of the first and the read end fd 0 of the second. Both run at the same time. When the buffer is full the writer waits, and when it's empty the reader waits.

```
tinybox cat a.txt ── fd 1 ──▶ [ pipe buffer ] ── fd 0 ──▶ wc -l
```

### EOF and SIGPIPE

When every write end of a pipe is closed and the buffer is drained, the reader's next read returns **0 bytes. That is end-of-file (EOF).** EOF isn't a character and it isn't a signal; it's a read that returns nothing. Pressing Ctrl-D in a terminal makes the pending read return whatever has been typed so far. If nothing has been typed, that's a 0-byte read, which is EOF.

It works the other way too. If the reader goes away while the writer is still writing, the writer's next write fails because the pipe is "broken", and by default the kernel sends the writer the **SIGPIPE** signal, which kills it. That's how `huge-output | head -n 1` finishes immediately instead of running forever.

## 4. Relevant Go concepts

### Errors are values

Functions that can fail return an `error` as their last result, and you check it explicitly:

```go
n, err := strconv.Atoi("42x")
if err != nil {
	fmt.Fprintln(os.Stderr, "bad number:", err)
	return
}
```

There's no `try`/`catch` and no rejected promise. `err` is an ordinary value, either `nil` or something with an `Error() string` method. Errors from the `os` package carry context: opening a missing file gives `open missing.txt: no such file or directory`, which names the operation, the path, and the cause.

The `fmt.Println(err); return` pattern from many tutorials (including Games With Go) fails a command-line tool twice: the message goes to stdout, and returning from `main` reports success.

### `*os.File`, `os.Open`, and `defer`

`os.Open(name)` opens a file for reading and returns an `*os.File` (a pointer to a `File`) and an error. `os.Stdin`, `os.Stdout`, and `os.Stderr` are `*os.File` values wrapping fds 0, 1, and 2. Close files you open once you're done with them, to give the descriptor back.

```go
func firstByte(name string) (byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close() // runs when firstByte returns, whichever way it returns
	// ... read from f ...
}
```

`defer` runs when the surrounding **function** returns, not when the surrounding block ends. That detail matters in this lesson's experiments.

### Interfaces: `io.Reader` and `io.Writer`

An **interface** is a set of method signatures. Any type that has those methods satisfies the interface automatically; there's no `implements` keyword. That's structural typing, as in TypeScript, but only for methods. Two tiny interfaces carry most of Go's I/O:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
```

`*os.File` has both methods, so a file is both a Reader and a Writer. So are network connections, in-memory buffers (`bytes.Buffer`), compressors, and hash functions. Code written against `io.Reader` works with all of them, which is how you'll test your applets without real files in Lesson 3.

**Read blocks.** There's no `await`. If no data is available yet, `Read` simply doesn't return until there is (or until EOF or an error). Coming from Node this can feel wrong. In Go, blocking is cheap because goroutines are cheap, and the runtime does the asynchronous part for you. You'll see how in P3.

The `Read` contract is subtle, and getting it wrong is a classic bug. From `go doc io.Reader`:

> Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered. … Callers should always process the n > 0 bytes returned before considering the error err.

So a single `Read` can return fewer bytes than you asked for, and it can return data *and* `io.EOF` from the same call. EOF arrives as the error value `io.EOF`, not as a special byte.

`io.Copy(dst, src)` runs the Read/Write loop for you until EOF. It returns the number of bytes copied and the first real error. A successful Copy returns `err == nil`, because reaching EOF is how a copy is supposed to end.

### Byte slices

`[]byte` is a slice of bytes, Go's type for raw data. `make([]byte, 4096)` allocates a 4096-byte buffer, and `buf[:n]` is the first `n` bytes of it: the same memory, not a copy.

### If you remember `io/ioutil`

It has been deprecated since Go 1.16. Its functions moved to `io` and `os`: `io.ReadAll`, `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, `io.Discard`, and so on.

## 5. The challenge

Use two terminals from now on: one in `projects/01-tinybox` for editing and building, and one in a throwaway directory for running.

```
# terminal 2: a scratch directory, with tinybox on PATH for this shell only
cd "$(mktemp -d)"
export PATH="$HOME/dev/limpet/projects/01-tinybox/bin:$PATH"
which tinybox
```

Add a `cat` applet:

```
tinybox cat [file...]
```

- Copy each named file to standard output, in order.
- With no file arguments, copy standard input to standard output. A file argument of `-` also means standard input.
- If a file can't be opened or read, print an error prefixed with `cat:` to standard error, then **carry on with the next file**.
- Exit with status 0 if everything succeeded, and 1 if anything failed.
- Close each file when you're done with it.

Use `io.Copy` for the copying. You'll write the loop yourself in an experiment.

## 6. Expected behavior

In terminal 2:

```
$ printf 'hello\n' > a.txt; printf 'world\n' > b.txt
$ tinybox cat a.txt b.txt
hello
world
$ tinybox cat a.txt missing.txt b.txt; echo "status=$?"
hello
cat: open missing.txt: no such file or directory
world
status=1
$ tinybox cat a.txt missing.txt b.txt 2>/dev/null
hello
world
$ printf 'middle\n' | tinybox cat a.txt - b.txt
hello
middle
world
$ tinybox cat < a.txt
hello
$ tinybox cat a.txt | tinybox cat | tinybox cat
hello
```

Error wording may differ (GNU cat says `cat: missing.txt: No such file or directory`). What can't differ: which stream the error goes to, that `b.txt` still gets printed, and the exit status.

## 7. Experiments

Predict first, every time.

1. **EOF from a keyboard.** Run `tinybox cat` with no arguments. Type a line and press Enter. Then type `abc` and press Ctrl-D *without* pressing Enter. What happens? Press Ctrl-D again. Explain both results using "a read returns whatever is available; a read that returns 0 bytes means EOF".

2. **Look at the fd table.** In terminal 2, run `tinybox cat > out.txt`; it waits for input. In terminal 1, run `lsof -p "$(pgrep -n tinybox)"`. Find fds 0, 1, and 2 in the FD column (`r`, `w`, and `u` mean read, write, and both). What does each point at, and who set that up? You may see entries your code never opened, such as a `KQUEUE` on macOS, which belongs to the Go runtime's I/O machinery. Note them; we'll meet them again in P2. Press Ctrl-D to finish.

3. **A file that becomes its own output.** Predict the contents of `a.txt` after `tinybox cat a.txt > a.txt`, then run it. Explain the exact order of events: which process opened `a.txt` for writing, when, and what that did to the file.

4. **Directories.** Run `tinybox cat .`. Predict which step fails, opening or reading, then check your error message.

5. **Broken pipes.** Create a large file with `seq 1 200000 > big.txt`, then run `tinybox cat big.txt | head -n 1`. In zsh, `echo $pipestatus` prints the status of every command in the last pipeline. Predict both numbers. Then read the SIGPIPE section of `go doc os/signal`: Go treats writes to fds 1 and 2 specially. Why might that be a sensible design?

6. **The defer trap.** Change `cat` to `defer f.Close()` inside the loop over files instead of closing each file explicitly. (If your structure makes that impossible, explain why your design is immune.) Create 30 small files and run cat with a lowered limit on open files. Use a subshell, so your terminal keeps its normal limit:
   ```
   for i in $(seq 1 30); do echo $i > f$i.txt; done
   ( ulimit -S -n 16; ulimit -H -n 16; tinybox cat f*.txt > /dev/null ); echo $?
   ```
   Predict whether it fails, and on which file. Then explain *exactly* which file it failed on. You'll need two separate facts: what `defer` does, and the order in which `f*.txt` expands. Put your close back afterwards.

7. **Write the loop yourself.** Replace `io.Copy` with your own loop: a `[]byte` buffer, `Read`, then `Write` whatever you got, until EOF, following the Read contract above exactly. Then measure how much the buffer size matters:
   ```
   dd if=/dev/zero of=big.bin bs=1048576 count=64     # 64 MiB of zeros
   hyperfine --warmup 1 'tinybox cat big.bin > /dev/null'
   ```
   Try buffers of 32 KiB, 4 KiB, 512 bytes, and finally 1 byte. For the 1-byte run, shrink the file to 1 MiB first, and predict why you need to. Every `Read` and every `Write` on a file is a system call: a trip into the kernel. What does that tell you about where the time goes? Afterwards, switch back to `io.Copy`, or keep your loop if you can argue why it's correct.

## 8. Understanding checks

1. In `tinybox cat a.txt | wc -l`, how many processes run, who creates the pipe, and what are fds 0 and 1 of each process connected to?
2. Describe a concrete pipeline where cat writing its error messages to stdout instead of stderr would silently corrupt data.
3. A `Read` returns `n = 5` and `err = io.EOF`. What must your code do with those 5 bytes, and why is a Reader allowed to do that?
4. Explain `tinybox cat a.txt > a.txt` step by step: which process opens `a.txt` for writing, when, and with which flags? (`man 2 open` has the flag names.)
5. Why does cat keep going after a missing file instead of stopping? Name a tool for which stopping at the first error would be the better design, and say why.

## 9. Completion criteria

- [ ] The expected behavior matches (wording aside).
- [ ] `gofmt -l .` is empty and `go vet ./...` is clean.
- [ ] Your hand-written copy loop is correct under the Read contract. We'll review it together.
- [ ] You can explain the fd table, redirection, pipes, EOF, and SIGPIPE in your own words, and you've answered the understanding checks.

## 10. Further exploration

- `man 2 read`, `man 2 write`, `man 2 open`, `man 2 dup2`. `dup2` is how a shell installs a file as fd 1, and limpet will do the same.
- [Errors are values](https://go.dev/blog/errors-are-values), Rob Pike on why Go handles errors this way.
- `cat -n` using `bufio.Scanner`. Then feed it a file with one 100 KB line. What happens, and why?
- `tee` using `io.MultiWriter`.
- Detect whether stdout is a terminal (`os.Stdout.Stat()` and `os.ModeCharDevice`). Why do `ls` and `grep --color=auto` behave differently when piped?
- Lowering only the soft limit (`ulimit -S -n 16`) doesn't stop your program from opening many files. Why not? Hint: the [Go 1.19 release notes](https://go.dev/doc/go1.19), under "Runtime".
- A teaser for P2: on Linux, `strace` shows every syscall your cat makes. Predict now what `io.Copy` turns into, and we'll check.
