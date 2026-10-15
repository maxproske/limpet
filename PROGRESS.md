# Progress

_Last updated: 2026-10-13_

## Now

- **Project:** P1 · tinybox (`projects/01-tinybox/`)
- **Lesson:** L1 · Hello, process (ready to start)
- **State:** course scaffolded and calibrated from the diagnostic answers. Go 1.27.2 is installed.

## Next actions

1. L1: `go mod init` in `projects/01-tinybox`, then the `true`, `false`, and `echo` applets, then the experiments.
2. Approve the scaffold's first commit when you're happy with it.

## Milestones

| Project | Status |
|---|---|
| P1 · tinybox: Unix tools in Go | in progress |
| P2 · Inside the process: /proc, syscalls, the fd table | not started |
| P3 · limpet run: spawn, wait, signal | not started |
| P4 · limpet init: a tiny init and supervisor | not started |
| P5 · Namespaces | not started |
| P6 · Root filesystems: chroot to pivot_root | not started |
| P7 · cgroups v2 | not started |
| P8 · limpet v0.1: a shell in a box | not started |
| P9 · Networking | not started |
| P10 · Hardening | not started |
| P11 · Speaking OCI | not started |
| P12 · Virtual machines from the outside | not started |
| P13 · Choose a VMM path | not started |
| P14 · limpet sandbox | not started |

## Current project: P1 lessons

| Lesson | Status |
|---|---|
| L1 · Hello, process | ready |
| L2 · Streams and file descriptors | ready (planned as two sessions) |
| L3 · Testable tools | sketch only (written after L2) |
| L4 · Who am I? (`env`, `probe`) | sketch only |
| L5 · One binary, many names | sketch only |

Implementation status only. What you understand is tracked in `LEARNER.md`.

## Parking lot

Open questions, requested detours, and ideas to revisit:

- After P7: review the "Containers Are Not Magic" series as a pull request, then compare with `docs/spoilers/containers-are-not-magic.md`.

## Decisions

- 2026-10-13: One Go module per project, created by you in each project's first lesson. Go 1.27 pinned via `mise.toml`.
- 2026-10-13: tinybox is a single multi-tool binary from L1, so applets accumulate in one place.
- 2026-10-13: The Linux lab VM is set up at the start of P2. Installing Lima then is approved. The M5 Pro supports the nested virtualization P12 needs.
- 2026-10-13: Sessions run 45–60 minutes, one a day. Lessons that need longer are planned as multiple sessions.
