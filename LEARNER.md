# Learner profile

Claude's working model of the learner. It's updated from evidence (explanations, predictions, code, debugging), never just because code compiles or tests pass.

_Last updated: 2026-10-13, after the diagnostic answers and a read of the learner's public Go repos._

## Background

| Area | What we know | Source |
|---|---|---|
| Primary development | Deep TypeScript: Next.js APIs, full-stack development | stated |
| Containers | **Docker Captain**: expert Docker user and community leader; little knowledge of container fundamentals (what the kernel does) | stated |
| Go, 2018 | Worked through ~20 chapters of Games With Go (`maxproske/games-with-go`): pointers, structs, methods, slices, linked lists, trees, graphs, goroutines with `sync`, packages, interfaces, SDL2 via cgo, noise, ASTs, genetic algorithms, a lexer. Also `gopl` and `idiomatic-go` (The Go Programming Language). GOPATH era, no modules. | repos |
| Go, 2026 | Restarted twice: *Learn Go with Tests* (`learn-go`, through ch. 5, structs/methods/interfaces, May) and *Learning Go* (`learn-go-oreilly`, through ch. 2, arrays and slices, July). Uses modules and basic tests. | repos |
| C | Second-year university C, a long time ago | stated |
| Linux | Has administered and hardened DigitalOcean droplets. Has never used `/proc`, `strace`, or cgroups. | stated |
| Other tooling | pnpm, Biome, uv (Python), mise, just; a modern CLI kit (rg, fd, ast-grep, jq, hyperfine) | global Claude config |
| Machine | Apple M5 Pro (arm64), 48 GB RAM, macOS 27.0.1, zsh. Docker Desktop 29.8; VS Code. Go 1.27.2 via mise (installed 2026-10-13). | observed |
| Pace | One session a day, 45–60 minutes | stated |

## Goals

- Become a capable Go and systems programmer.
- Understand containers and VMs by building small versions of their mechanisms.
- Independence: read man pages and kernel docs, design solutions, and debug without AI writing the code.

## Learning preferences (stated)

- Project-based and hands-on, in the spirit of Handmade Hero and Games With Go.
- Guided discovery: attempt first, progressive hints, solutions only on explicit request.
- Wants to be asked to predict, explain, and experiment.
- Wants rigorous feedback: no praise for wrong answers, oversimplifications flagged.
- Short feedback loops; working programs over abstract exercises.

## Teaching notes

- **Go syntax is a refresher, not new.** Skim the basics; spend session time on systems concepts and on what changed since 2018. Flag changes as they become relevant (all verified in the release notes): `errors.Is`/`errors.As`/`%w` (1.13), `io/ioutil` deprecated (1.16), generics (1.18), `slices`/`maps`/`log/slog` and `min`/`max` (1.21), per-iteration loop variables (1.22), `os.Root` (1.24), `testing/synctest`, `sync.WaitGroup.Go`, and container-aware `GOMAXPROCS` (1.25).
- **2018 habits to watch for:** `GetX` getters (28 in games-with-go; idiomatic Go drops the `Get`), and `fmt.Println(err); return` error handling (errors belong on stderr, with context and a non-zero exit). In tests, Go convention is `got`/`want` rather than `actual`/`expected`.
- **Writes explanatory notes in code comments.** Some are confident oversimplifications, such as "structs are Go's only user-defined type" (named function, map, and slice types can have methods too, which tinybox will use). Review comments as carefully as code. The habit is a strength for concept notes: start them from these explanations.
- **Already motivated by limpet:** May 2026 notes mention mapping the OCI runtime spec onto Go structs.
- **Two book-paced restarts stopped early.** Keep every session moving something real forward, and keep drill to a minimum.
- **Bridge from what's deep:** TypeScript and Node for Go concepts (errors vs. exceptions, blocking I/O vs. `await`, structural interfaces); Docker behavior for systems concepts (`Exited (137)`, `--init`, `docker logs` streams, `-it`). Docker expertise is also a trap: knowing which flag to pass isn't knowing what the kernel does with it.
- **Session shape:** 45–60 minutes. Each session should end at a natural stopping point, with a short recap and a quick wrap-up.

## Understanding tracker

Levels: **seen** (met in a lesson) → **used** (in working code, possibly with help) → **explains** (explained it correctly, in your own words) → **applies** (used it correctly in a new situation, unprompted). Mark **shaky** with the misconception when one appears.

| Concept | Level | Evidence (date: what showed it) |
|---|---|---|
| Pipelines create one process per command | explains | 2026-10-13: calibration, "cat and grep are two processes" (correct) |
| Who performs redirection (`> out.txt`) | **shaky** | 2026-10-13: calibration, "nothing opens out.txt". In fact the shell opens it in the child before exec. Target: L2. |
| Who enforces container resource limits | **shaky** | 2026-10-13: calibration, "maybe the Docker engine keeps a ledger JSON". Docker does store settings as JSON, but the kernel enforces limits through cgroups. Target: P7. |
| Go basics: slices, structs, methods, interfaces, simple tests | used (2018, 2026) | repos; not yet demonstrated in this course |

## Strengths

- Broad programming experience and a substantial earlier Go project, so the language should come back quickly.
- Hands-on Linux server administration (droplet hardening).
- Writes down explanations while learning.

## Gaps and watch-list

- The kernel's process and file model (redirection, fds), first seen in the calibration answers. L2 targets it.
- Possible oversimplifications in self-written notes. Check them.
