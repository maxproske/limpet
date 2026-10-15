# limpet

A personal, long-running course in Go and Linux systems programming, organized around one destination: **limpet**, a small Linux sandbox runtime written from first principles in Go. It grows from tiny Unix tools into a process supervisor, then a container runtime that can run a shell, and perhaps a microVM-backed sandbox for development tools and coding agents.

The point isn't the finished runtime. The point is becoming someone who could build one: someone who knows what the kernel is doing, writes idiomatic Go, and can go from a man page to a working design without help.

Inspired by Casey Muratori's [Handmade Hero](https://handmadehero.org/) and Jack Mott's [Games With Go](https://gameswithgo.org/): build real things from the ground up, understand why they work, and break them on purpose.

## How it works

- **You write the code.** Claude, running in Claude Code in this directory, is the instructor and reviewer. It explains, asks you to predict and explain, reviews your work, and gives hints when you ask. It shows a full solution only when you explicitly request one.
- **The roadmap is written up front; lessons are written just in time.** [`ROADMAP.md`](ROADMAP.md) lays out the whole arc. Detailed lessons exist only for the current project, and each one is written from how the previous ones went.
- **Understanding is tracked separately from progress.** [`PROGRESS.md`](PROGRESS.md) records what's built. [`LEARNER.md`](LEARNER.md) records what you can explain and apply, with evidence.

## A session

1. Open Claude Code in this directory and say something like "let's continue". Claude reads the progress files, looks at your code, recaps where you stopped, and proposes the next step.
2. Work through the lesson in small steps: predict, run, observe, explain.
3. Stuck? Ask for a **hint**, then a **stronger hint** if you need one. Want the answer? Say **"show me the solution"**.
4. Finished an exercise? Ask for a **review**.
5. Say **"wrap up"** to finish. Claude writes a journal entry and updates the progress file, the learner profile, and the concept notes.

## Layout

```
CLAUDE.md          how Claude teaches here (loaded every session)
PROMPT.md          the founding brief
ROADMAP.md         the curriculum: parts, projects, decision points, detours
PROGRESS.md        where you are and what's next
LEARNER.md         background, preferences, and an evidence-based understanding tracker
projects/NN-name/  one Go module per project: outline, lessons/, and your code
docs/concepts/     short notes on concepts you've actually worked through
docs/resources.md  curated references, with spoilers marked
experiments/       disposable investigations
journal/           dated session notes
mise.toml          the pinned Go toolchain
```

## Where things run

- **The macOS host** runs Go itself and Project 1. Nothing privileged ever runs on the host.
- **A disposable Linux VM**, set up at the start of Project 2, runs everything Linux-specific: `/proc` and `strace` at first, later namespaces, mounts, cgroups, and KVM. Root is used only inside that VM.
- **Docker Desktop** is used only for unprivileged side-by-side comparisons, such as running your `probe` in an ordinary container.

limpet is a teaching prototype. At no stage is it a secure sandbox, and nothing in this repo claims otherwise.

## Toolchain

Go is pinned in [`mise.toml`](mise.toml). Run `mise install` in this directory, then `go version`.
