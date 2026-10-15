# Lesson 1 · Hello, process

> P1 tinybox · runs on macOS · prerequisites: none

## 1. Objective

By the end of this lesson you will:

- have a working Go toolchain and a Go module for tinybox;
- have written `tinybox` with three applets: `true`, `false`, and `echo`;
- be able to explain what happens between pressing Enter and seeing `$?`: who turns your command line into arguments, who runs your program, and how your program reports success or failure.

## 2. Why it matters

limpet's job is to start programs with carefully controlled inputs (arguments, environment, open files) and to report how they ended. Before you control those things for *other* programs, look at them from inside your own. Everything here comes back in P3, when you write code that does the shell's job of starting and waiting for processes.

## 3. The mental model

### Program vs. process

A **program** is a file. `bin/tinybox` on disk is just bytes: machine code plus metadata telling the OS how to load it. A **process** is a running instance of a program, with its own memory, its own open files, its own process ID (pid), and a parent that started it. Run tinybox three times at once and you get three processes from one program.

You've been running processes for years, but in Node the program the kernel runs is `node`, and your JavaScript is data that `node` interprets. That's why `process.argv[0]` is the path to `node` and `process.argv[1]` is your script. A Go program compiles to a native executable, so here the kernel runs *your* program directly.

### What happens when you type `./bin/tinybox echo "hello world"`

```
zsh reads the line and parses it:
  expands $VARIABLES and *globs, removes quotes, splits words
  → argv = ["./bin/tinybox", "echo", "hello world"]
        │
        ▼
zsh creates a child process (fork), then waits for it
        │
        ▼
the child asks the kernel to run a different program in itself:
  execve("./bin/tinybox", argv, environment)
  the kernel loads tinybox into that process, places argv and the
  environment in its memory, and jumps to the program's entry point
  (on macOS by way of the dynamic loader, dyld; see Experiment 7)
        │
        ▼
the Go runtime initializes, then calls your main()
  your code reads os.Args and writes to standard output
        │
        ▼
the process exits with a status number;
the kernel keeps that number until the parent collects it
        │
        ▼
zsh's wait returns and zsh stores the status in $?
```

Three things to notice:

1. **Your program never sees the command line.** It receives an array of strings. Quoting, `$VARIABLES`, and `*` globs are all handled by the shell *before your program exists*; the kernel just passes the array along.
2. **Starting a program takes two steps in Unix:** create a process (`fork`), then replace the program running in it (`execve`). Container runtimes live in the gap between those steps, where they change things about the child (its view of the system, its limits) before the real program starts. Some launchers use faster variants (`vfork`, or `posix_spawn`, which does both steps in one call), but the model is the same. You'll write this yourself in P3 and discover why Go doesn't offer a plain `fork`.
3. **The exit status is tiny.** The parent sees only the low 8 bits of the number you exit with, so 0–255. By convention, 0 means success and anything else means failure. If a process is *killed by a signal* instead of exiting, the parent learns that too, and shells report it as 128 + the signal number. That's the `Exited (137)` you've seen in `docker ps -a`: 137 = 128 + 9, and signal 9 is SIGKILL, which is what the OOM killer and a forced stop use.

### Standard output and standard error

Every process starts with three open channels: **standard input** (stdin), **standard output** (stdout), and **standard error** (stderr). Normal output goes to stdout. Diagnostics such as usage messages and errors go to stderr, so they still reach a human when stdout is redirected into a file or piped into another program. Lesson 2 is about what these channels really are. For now, keep the two outputs separate.

## 4. Relevant Go concepts

You've written all of this before, in games-with-go and *Learn Go with Tests*. Treat the section as a refresher: skim it, and slow down only where something surprises you. The systems ideas in section 3 are where this lesson's time should go.

### Modules and packages

A **module** is a unit of Go code defined by a `go.mod` file at its root, which records the module's path and the Go version it targets (and later its dependencies). It plays roughly the role of a `package.json`. A **package** is a directory of `.go` files compiled together, and every file in it starts with `package <name>`. There's no per-file import/export: identifiers starting with an uppercase letter are visible outside the package, and lowercase ones aren't.

`package main` is special. It builds into an executable, and its `func main()` is where your code starts:

```go
package main

import "fmt"

func main() {
	fmt.Println("hi")
}
```

### `os.Args` and slices

`os.Args` holds the arguments, starting with the program name, so `os.Args[0]` is the program as it was invoked. Its type is `[]string`, a **slice**: Go's growable view onto an array.

```go
words := []string{"alpha", "beta", "gamma"}
len(words)   // 3
words[0]     // "alpha"
words[1:]    // []string{"beta", "gamma"}
words[3:]    // []string{} (slicing *at* the length is allowed; the result is empty)
words[3]     // panics: index out of range
```

Indexing past the end doesn't give you `undefined`. It **panics**: the program stops with an error message and a stack trace. Go reserves panics for bugs. Expected failures are ordinary return values, which you'll meet in Lesson 2.

### `switch`

```go
switch day {
case "sat", "sun":
	fmt.Println("weekend")
case "":
	fmt.Println("no day given")
default:
	fmt.Println("weekday")
}
```

Cases don't fall through to the next one (unlike in C and JavaScript), so there's no `break`, and one case can list several values.

### Printing and joining

- `fmt.Println(a, b)` prints its operands separated by spaces, then a newline.
- `fmt.Print(s)` prints without adding a newline. With several operands, it adds spaces only between operands when neither is a string, so read `go doc fmt.Print` before relying on it.
- `fmt.Fprintln(os.Stderr, "msg")` prints to standard error instead of standard output.
- `strings.Join([]string{"a", "b", "c"}, "-")` returns `"a-b-c"`.

### Exiting

- Returning from `main` exits with status 0.
- `os.Exit(code)` exits immediately with `code`. Read its documentation now (`go doc os.Exit`); there's an important sentence about `defer` that you'll test in the experiments. `defer` schedules a call to run when the surrounding function returns. It's Go's closest equivalent to `finally`.

### The toolchain

- `go run . args…` builds the package in the current directory and runs the result. Since Go 1.24 the build is cached, so a second run with no changes starts quickly.
- `go build -o bin/ .` builds it into `bin/tinybox`.
- `gofmt -l .` lists files that aren't in standard format, and `go fmt ./...` fixes them. Go has exactly one formatting style.
- `go vet ./...` reports suspicious code that compiles anyway.
- `go doc os.Args`, `go doc strings.Join`, and so on show documentation in the terminal.

## 5. The challenge

### Step 0 · The toolchain

Go 1.27.2 is already installed through mise, which reads the version pinned in `mise.toml` (three lines; worth a look). From anywhere inside the repository:

```
go version          # expect: go version go1.27.2 darwin/arm64
```

Outside this repository, `go` reports that no version is set, because mise only activates Go here. If you edit Go in VS Code, install the official Go extension, which sets up `gopls`, the Go language server.

### Step 1 · The module

```
cd projects/01-tinybox
go mod init limpet/tinybox
```

Open the `go.mod` it created. Can you say what each line means? (We'll discuss it.)

Published modules use paths like `github.com/you/project` so that `go get` can find them. Nobody will ever download this module, so a short local path is fine.

### Step 2 · tinybox

Create `main.go` in `projects/01-tinybox`, in `package main`, so that:

- `tinybox true` prints nothing and exits with status 0.
- `tinybox false` prints nothing and exits with status 1.
- `tinybox echo [-n] [words…]` prints its words separated by single spaces, followed by a newline. If the first word is exactly `-n`, that word isn't printed, and neither is the final newline.
- `tinybox` with no applet name prints a short usage message to **stderr** and exits with status 2.
- `tinybox <unknown>` prints an error naming the unknown applet to **stderr** and exits with status 2.

Exit status 2 for "you used the command wrong" is a widespread convention, which Go's own `flag` package follows.

Keep everything in one file and write the smallest thing that works. In Lesson 3, tests will give us a reason to restructure.

## 6. Expected behavior

Your usage and error wording can differ. The streams and the statuses can't.

```
$ go build -o bin/ .
$ ./bin/tinybox echo hello world
hello world
$ ./bin/tinybox echo -n no newline; echo "<- the next output starts here"
no newline<- the next output starts here
$ ./bin/tinybox echo

$ ./bin/tinybox true; echo $?
0
$ ./bin/tinybox false; echo $?
1
$ ./bin/tinybox; echo $?
usage: tinybox <applet> [args...]
2
$ ./bin/tinybox frobnicate; echo $?
tinybox: unknown applet "frobnicate"
2
$ ./bin/tinybox frobnicate 2>/dev/null; echo $?
2
```

The last command shows that your error went to stderr: `2>/dev/null` discards stderr, and nothing else is printed.

If you run `./bin/tinybox echo -n hi` on its own, zsh shows a highlighted `%` after `hi`. That's zsh marking output that didn't end in a newline. Your program didn't print it.

## 7. Experiments

Before each experiment, write down your prediction, even if it's a guess. A wrong prediction that you then explain teaches more than a right one you can't explain.

1. **Who splits the words?** Predict the output of each command, then run it:
   ```
   ./bin/tinybox echo hello,      world
   ./bin/tinybox echo "hello,      world"
   ./bin/tinybox echo '$HOME' "$HOME"
   ./bin/tinybox echo *.go
   ```
   Then temporarily change `echo` to print each argument on its own line inside brackets, like `[hello,]`, and run them again. Which of these transformations did your program perform? Change it back afterwards.

2. **What is `os.Args[0]`?** Print it to stderr, so it doesn't pollute echo's output, and compare `./bin/tinybox echo hi`, `bin/tinybox echo hi`, `"$PWD/bin/tinybox" echo hi`, and `go run . echo hi`. Where does `go run` put the program it builds?

3. **Exit statuses, up close.** Add a temporary applet that exits with status 3.
   - Predict `$?` after `./bin/tinybox <it>` and after `go run . <it>`, then explain the difference. How many processes were involved in the second case, and whose status did zsh see?
   - Change the status to 256, then to -1. Predict `$?` before each run.
   - Have that applet `defer` a print to stderr before calling `os.Exit`. Does the print happen? Now make an applet that defers a print and then *returns* from `main`. Why might this difference matter to limpet, which will need to clean up (unmount filesystems, delete cgroups) before it exits?

4. **Break it.** Remove your check for a missing applet name, rebuild, and run `./bin/tinybox` with no arguments. Read the whole panic message. What's the exit status, and which stream did the message go to? (Check with `2>/dev/null`.) Restore the check.

5. **Killed vs. exited.** Run `sleep 30`, press Ctrl-C, then `echo $?`. How is "killed by a signal" different from "exited with a status"? Which signal was it?

6. **Not every command is a process.** Run `type echo cd ls` in zsh. Why would a shell implement some commands itself? Which of these three *must* be built into the shell, and why couldn't it work as a separate program?

7. **What did you build?** Run `file bin/tinybox`, `ls -lh bin/tinybox`, and `otool -L bin/tinybox`. What kind of file is it, why is it so large for such a small program, and what is it linked against? Keep your answers; on Linux in P2 they'll differ in an interesting way.

## 8. Understanding checks

Answer in your own words, in the chat. Partial answers and "I'm not sure" are useful, because they tell us what to do next.

1. For `./bin/tinybox echo "$HOME"/*.txt`, list everything that happens to that text before your program starts. Which parts does zsh do, which does the kernel do, and which does your program do?
2. A CI script runs `go run ./cmd/deploy` and treats `$?` of 3 as "already deployed, skip". Will that work? Explain why or why not, and what you'd change.
3. Why should a usage error go to stderr, even when the user sees the same terminal either way? Describe a concrete situation where getting this wrong causes a bug.
4. limpet will start programs for you and report how they ended. Based on this lesson, what are the two broad ways a process can end, and what does the parent learn in each case?

## 9. Completion criteria

- [ ] The expected behavior above matches (wording aside).
- [ ] `gofmt -l .` prints nothing and `go vet ./...` is clean.
- [ ] You've run the experiments and can explain each result, especially the ones you predicted wrong.
- [ ] You've answered the understanding checks and we've discussed them.

## 10. Further exploration

- `man 2 execve` and `man 2 wait` on your Mac. These are the BSD versions; we'll read the Linux ones in P2.
- The [Go Tour](https://go.dev/tour/), "Basics" section, for more on variables, types, and functions.
- Run `docker run --rm alpine sh -c 'exit 3'; echo $?`. Does `docker run` behave like `go run` here? Then look up what Docker means by exit statuses 125, 126, and 127.
- `go doc os.Exit` recommends statuses in [0, 125]. What do 126, 127, and 128+ mean to a shell?
- Read the source with `go doc -src os.Exit`. What does it do before asking the kernel to end the process?
- Add `-e` to `echo` (escape sequences such as `\n` and `\t`). Different `echo` implementations disagree about this, which is why portable scripts use `printf`.
