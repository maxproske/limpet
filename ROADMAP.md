# Roadmap

**Destination:** `limpet`, a small Linux sandbox runtime in Go, built from first principles. It starts as a process runner, becomes a container runtime that can run a shell, and may end as a microVM-backed sandbox for dev tools and coding agents.

This is a direction, not a spec. At the end of every project we ask whether the next step is still the most interesting and educational one. The aim isn't to recreate Docker, runc, or Firecracker. It's to build small versions of their mechanisms until those mechanisms stop being magic.

Status lives in `PROGRESS.md`. Lessons are written just in time, inside each project's folder.

## The arc

```
Part I    Go and the Unix process model     P1 tinybox → P2 inside the process → P3 limpet run → P4 limpet init
Part II   Isolation from first principles   P5 namespaces → P6 root filesystems → P7 cgroups v2 → P8 limpet v0.1: a shell in a box
Part III  Toward a real runtime             P9 networking → P10 hardening → P11 OCI
Part IV   Virtual machines                  P12 VMs from the outside → P13 choose a VMM path
Part V    The destination                   P14 limpet sandbox
```

**Where things run.** Only P1 runs on the macOS host. From P2 on, work happens in a disposable Linux VM, and from P5 on, experiments need root *inside that VM*. Privileged work never touches the host.

**Sizes** (S < M < L < XL) hint at relative effort. They aren't a schedule.

## From Docker flags to kernel mechanisms

You already know what containers *do* from the outside. Each part of this course opens up something you've used:

| You've used… | What's underneath | Built in |
|---|---|---|
| `Exited (137)` in `docker ps -a` | exit statuses; 128 + signal number (9 = SIGKILL) | P1, P3 |
| `docker logs` keeping stdout and stderr apart (until you add `-t`) | file descriptors and pipes | P1, P3 |
| `docker stop` hanging for 10 s on some images; `--init` fixing it | PID 1 signal semantics, zombie reaping | P4 |
| `--hostname`, `--pid=host`, `--network none`, `--ipc`, userns-remap | namespaces | P5 |
| `docker exec` | joining existing namespaces with `setns` | P5 (extension) |
| image layers, `--read-only`, `--tmpfs`, bind mounts, volumes | mount namespaces, `pivot_root`, overlayfs | P6 |
| `--memory`, `--cpus`, `--pids-limit`, OOM kills | cgroups v2 | P7 |
| `-it` | pseudo-terminals, sessions, job control | P8 |
| `-p 8080:80`, the default bridge network | veth pairs, bridges, NAT | P9 |
| `--cap-drop`, `--security-opt seccomp=…`, `no-new-privileges`, `--privileged`, `--user` | capabilities, seccomp-BPF, privilege dropping | P10 |
| `docker` → containerd → runc, `--runtime` | the OCI runtime spec and lifecycle | P11 |
| Docker Desktop's Linux VM, Docker Sandboxes | hypervisors, KVM, microVMs | P12–P14 |

A recurring experiment: run the same `probe` under `docker run` and under limpet, and explain every difference.

## Through-lines

- **`probe`, the test subject.** In P1 you write `probe`, which prints a process's view of the world: pid, user and group ids, working directory, hostname, environment. Every later project teaches it something new (namespaces, mounts, cgroups, capabilities, seccomp) and then runs it *inside* the new sandbox. If probe can't see a difference, the isolation didn't happen.
- **Predict → observe → explain.** The observation tools grow with the course: `lsof` (P1) → `strace` and `/proc` (P2) → your own `lps` (P2) → `unshare` and `nsenter` (P5) → `findmnt` (P6) → the cgroup filesystem (P7) → `ip`, `nft`, `tcpdump` (P9) → QEMU's monitor (P12).
- **An honest leak list.** From P5 on, `docs/threat-model.md` records what limpet isolates, what it doesn't, and how you know. That's how "teaching prototype" stays distinct from "secure sandbox".
- **One evolving codebase.** From P3 on, limpet is a single codebase that grows milestone by milestone, with a git tag at each one, so its history reads like Handmade Hero's day-by-day source. Where it lives gets decided at P3.
- **Dependencies are earned.** Standard library first. `golang.org/x/sys/unix` when `syscall` runs out (P2/P5). A third-party library only after you've built a minimal version yourself, so you know what you're delegating (netlink in P9, OCI types in P11).

---

## Part I: Go and the Unix process model

### P1 · tinybox: Unix tools in Go · M · macOS
- **Build:** `tinybox`, one Go binary containing small Unix tools: `true`, `false`, `echo`, `cat`, `env`, and `probe`. You run them as `tinybox <tool> args…`, and by the end through symlinks, BusyBox-style, selected by `argv[0]`.
- **Why:** You learn Go's core on problems where every line touches something limpet depends on: argv, exit status, the standard streams, the environment, process identity. The project also produces two things you'll reuse: `probe`, and a static binary that can be the only program in a root filesystem.
- **Learn.** *Go:* modules and packages, `main`, slices, strings, maps, structs, `switch`, errors as values, `defer`, `io.Reader`/`io.Writer`, `bufio`, `flag.FlagSet`, table-driven tests, `encoding/json`, `gofmt`/`go vet`, cross-compilation. *Systems:* program vs. process, how the shell builds argv, exit status (8 bits, 128 + n), file descriptors 0/1/2, redirection and pipes, EOF, SIGPIPE, the environment, pid/ppid/uid/gid/cwd, static vs. dynamic linking.
- **Prerequisites:** none.
- **Observable result:** `tinybox cat notes.txt missing.txt | tinybox cat` behaves like the real tools, including stderr, exit statuses, and broken pipes. A Linux build of tinybox runs in a `FROM scratch` container containing nothing else, and `tinybox probe` there reports pid 1.
- **Done when:** every tool handles pipes, redirection, bad input, and exit statuses correctly; each tool is tested through injected stdin/stdout/stderr; you can explain argv construction, exit statuses, and the fd table without notes.
- **Extensions:** `wc`, `head`, `tee` (`io.MultiWriter`), `yes` (make it fast), `cat -n` (and `bufio.Scanner`'s line-length limit), `seq`.
- **Leads to:** P2 runs tinybox under `strace`; P3 turns `env` into a real exec; P6 makes tinybox the userland of a root filesystem.

### P2 · Inside the process: /proc, syscalls, the fd table · M · Linux VM (unprivileged)
- **Build:** your Linux lab, a disposable VM you can recreate with one command. Then `lps`, a `ps`/`pstree` built by reading `/proc`, and Linux-only additions to `probe` (namespaces, cgroup, capabilities, mounts, open fds).
- **Why:** The kernel publishes process state as files under `/proc`, and every later milestone is observed through it. `strace` shows the syscall boundary under all of Go's abstractions, including the Go runtime's own threads.
- **Learn.** *Go:* parsing with `strconv`/`strings`/`bufio`, `os.ReadDir`, error wrapping (`%w`, `errors.Is`/`errors.As`), reading a system that changes while you read it, build constraints (`_linux.go`), cross-compiling with `GOOS`/`GOARCH`, `syscall` vs. `golang.org/x/sys/unix` (your first dependency decision). *Systems:* user vs. kernel mode, system calls, `strace`, process states, the process tree, the fd table vs. open file descriptions, `O_CLOEXEC`, signals observed from outside (`kill -STOP`/`-CONT`), resource limits.
- **Prerequisites:** P1.
- **Observable result:** `lps --tree` on the VM, and an strace of your own `tinybox cat` that you can annotate line by line.
- **Done when:** lps is correct and survives processes vanishing mid-scan; you can explain each group of syscalls in an strace of your program, including the ones the Go runtime makes on its own; probe's Linux-only parts build only on Linux while the macOS build keeps working.
- **Extensions:** an `lsof`-lite (who has this file open?), a `top`-lite (CPU% from `/proc/[pid]/stat` deltas), comparing an strace of GNU `cat` with yours.
- **Leads to:** P3 (you can watch the processes you spawn), P5 (namespaces show up as `/proc/[pid]/ns/*`).

### P3 · limpet run: spawn, wait, signal · M · Linux VM
- **Build:** the first limpet command, `limpet run [--timeout D] [--env K=V] -- cmd args…`. It starts a child with controlled argv, environment, and stdio; forwards signals; enforces a timeout by killing the whole process tree; and reports exactly how the child ended.
- **Why:** Every container is a process that something started, watched, and cleaned up after. This is the skeleton every later milestone hangs on.
- **Learn.** *Go:* `syscall.Exec` vs. `os.StartProcess` vs. `os/exec`, `exec.Cmd`, `syscall.SysProcAttr`, goroutines and channels (your first concurrency), `os/signal`, `select`, `context` and timeouts, decoding wait statuses, CLI subcommands. *Systems:* fork, exec, and wait, and why Go gives you no plain `fork()`; exec replacing the program while keeping the pid; zombies and reaping; exiting vs. dying by signal; signal dispositions across exec; process groups and sessions; `PR_SET_PDEATHSIG`; fd inheritance (`ExtraFiles`, `O_CLOEXEC`); `PATH` lookup; resource usage.
- **Prerequisites:** P2.
- **Observable result:** `limpet run --timeout 2s -- sh -c 'sleep 100 & sleep 100'` kills both sleeps and says so; `limpet run -- false` reports exit status 1; Ctrl-C reaches the child.
- **Done when:** exit codes and signals are reported correctly; timeouts leave no stray processes (checked with lps); tests pass under `-race`; you can explain why Go can't simply `fork()`.
- **Extensions:** waiting via pidfds, `--uid`/`--gid` via `SysProcAttr.Credential`, resource-usage reports.
- **Leads to:** P4 supervises what P3 spawns. P5 adds `Cloneflags` to the same `SysProcAttr`. This is a natural moment for detour D1, a tiny shell.

### P4 · limpet init: a tiny init and supervisor · M–L · Linux VM
- **Build:** `limpet init`, a minimal PID 1. It starts children, reaps zombies (including orphans it inherits), forwards signals, restarts failed children with backoff, shuts down gracefully (SIGTERM, a grace period, then SIGKILL), and exits with the main child's status. You test it as a subreaper first; it becomes a real PID 1 in P5 and P12.
- **Why:** The first process in a container, and in a VM, has special duties imposed by the kernel. That's why `docker run --init` exists, and why `docker stop` can hang for its full 10 seconds when PID 1 doesn't deal with SIGTERM itself. Plain `node` as PID 1 is the classic case, for a subtle reason you'll be able to explain by the end. This is also where Go concurrency stops being a toy.
- **Learn.** *Go:* goroutine lifecycles, channel patterns (fan-in of child events), `select` with timers, `sync` vs. channels, cancellation trees with `context`, the race detector, `log/slog`, testing concurrent code without sleeps. *Systems:* PID 1 semantics (signals without handlers, orphan reparenting), `SIGCHLD` and reaping loops, `PR_SET_CHILD_SUBREAPER`, graceful-shutdown protocols, restart storms and backoff.
- **Prerequisites:** P3.
- **Observable result:** under a deliberately badly behaved workload (double-forking orphans, children that ignore SIGTERM), `lps` shows no zombies and shutdown finishes within the grace period.
- **Done when:** that holds under `go test -race`; restart policies match their written spec; you can explain why the kernel doesn't deliver SIGTERM to a namespace's PID 1 unless it has installed a handler.
- **Extensions:** supervising several services from a config file, readiness checks, handing listening sockets to children (systemd-style socket activation).
- **Leads to:** P5 (PID 1 of a PID namespace), P8, P12 (PID 1 of a VM).

## Part II: Isolation from first principles

### P5 · Namespaces: changing what a process can see · L · Linux VM (root)
- **Build:** limpet gains `--ns` (uts, pid, mount, ipc, net, user, cgroup). The child starts in new namespaces, then a re-executed limpet sets up the child's view (hostname, a private `/proc`) before exec'ing the target. Optionally, a rootless mode through a user namespace.
- **Why:** Namespaces are the "what can I see" half of a container. You'll find that a container is a process with a different view of the system, and you'll find exactly where that view leaks.
- **Learn.** *Go:* `Cloneflags`, `Unshareflags`, `UidMappings`/`GidMappings`, the re-exec pattern (`/proc/self/exe`) and why Go needs it, `runtime.LockOSThread`, `golang.org/x/sys/unix`. *Systems:* the eight namespace types; `clone`, `unshare`, and `setns`; `/proc/[pid]/ns`; PID namespaces (PID 1 again, and the `/proc` remount trap); user namespaces and uid/gid maps; what namespaces don't hide. The experiments start with `unshare(1)` and `nsenter(1)` before any Go.
- **Prerequisites:** P4.
- **Observable result:** probe inside the sandbox shows its own hostname, pid 1, its own `/proc`, and, in rootless mode, uid 0 inside that maps to your unprivileged uid outside.
- **Done when:** you can predict probe's output for any combination of namespaces before running it, and you can demonstrate something namespaces don't isolate.
- **Extensions:** `limpet exec` into a running sandbox (an nsenter-lite), time namespaces.
- **Leads to:** P6 (the mount namespace is the stage for a new root), P7 (cgroup namespace), P9 (network namespace).

### P6 · Root filesystems: from chroot to pivot_root · L · Linux VM (root)
- **Build:** limpet gains `--rootfs DIR`. It makes mount propagation private, bind-mounts the new root, mounts `/proc`, a minimal `/dev`, and `/sys`, then calls `pivot_root` and detaches the old root. It starts with a root filesystem containing only tinybox, then moves on to BusyBox or Alpine. Optional overlayfs keeps the base image pristine.
- **Why:** This is the "what files exist" half of a container. `chroot` alone is famously escapable, and understanding why teaches what real filesystem isolation requires.
- **Learn.** *Go:* mount-related syscalls, cleanup that runs on every path (`defer`, `errors.Join`), path-traversal-safe joins, static builds (`CGO_ENABLED=0`) and why they matter here. *Systems:* the mount table; mount propagation (shared, private, slave, and why systemd makes `/` shared); bind mounts; pseudo-filesystems (`proc`, `sysfs`, `tmpfs`, `devtmpfs`); device nodes; a chroot escape demonstrated in the VM; `pivot_root`'s rules; overlayfs layers.
- **Prerequisites:** P5.
- **Observable result:** inside the sandbox, `ls /` shows only the root filesystem, and the VM's own mount table is unchanged afterwards.
- **Done when:** you can demonstrate a chroot escape, explain why `pivot_root` plus unmounting the old root defeats that particular escape, and list what a root process inside could still do. That list feeds P10.
- **Extensions:** a read-only root, overlayfs layers, a `tmpfs` `/tmp`, fetching a minirootfs and verifying its checksum.
- **Leads to:** P8; P11 (an OCI bundle is a config plus a root filesystem); P12 (an initramfs is a root filesystem too).

### P7 · cgroups v2: limiting what a process can use · M · Linux VM (root)
- **Build:** limpet gains `--memory`, `--cpus`, and `--pids`. Each sandbox gets its own cgroup, the child is placed in it atomically at creation, usage and OOM events are reported, and everything is removed on exit.
- **Why:** This is the "how much can I use" half of a container. Without limits, one sandboxed fork bomb or memory leak takes down the host.
- **Learn.** *Go:* the kernel's file-based APIs (writes that fail in instructive ways), `SysProcAttr.UseCgroupFD`/`CgroupFD`, integration tests that need root and skip cleanly otherwise. *Systems:* the unified hierarchy, controllers and `cgroup.subtree_control`, the no-internal-processes rule, `memory.max`/`memory.high` and the OOM killer, `cpu.max`, `pids.max`, delegation and systemd's ownership of the tree, cgroup namespaces, cgroups vs. rlimits.
- **Prerequisites:** P5 (P6 helps).
- **Observable result:** a fork bomb and a memory hog are contained and reported. The sandbox's cgroup is visible while it runs and gone afterwards.
- **Done when:** that holds, no exit path leaks a cgroup, and you can explain how the kernel enforces each limit.
- **Extensions:** pressure-stall (PSI) metrics, freezing a sandbox with `cgroup.freeze`, `io.max`, a live `limpet stats`. Recommended: review the "Containers Are Not Magic" tutorial series (see `docs/resources.md`) as if it were a pull request, then compare your findings with its errata.
- **Leads to:** P8; P13 (VMM processes get cgroups too, which is part of what Firecracker's jailer sets up).

### P8 · limpet v0.1: a shell in a box · L · Linux VM (root) · first major milestone
- **Build:** P3–P7 integrated. `limpet run --rootfs ./alpine --hostname box --memory 256M --pids 128 -- /bin/sh` gives you an interactive shell with working job control: a PTY, raw-mode relay, window resizing, `limpet init` as PID 1, and clean teardown on every exit path. Also the first version of `docs/threat-model.md`.
- **Why:** This is the first point where limpet is recognizably a container runtime. Integration is a skill of its own: ordering setup steps, reporting errors across a process boundary, cleaning up after failure.
- **Learn.** *Go:* package design (`internal/`), interfaces only at real seams, passing setup errors from child to parent over a pipe, `golang.org/x/term`, end-to-end tests. *Systems:* terminals and PTYs, sessions and controlling terminals (`setsid`, `TIOCSCTTY`), job-control signals, the line discipline and raw mode, `SIGWINCH`, ordering constraints (user namespace first, cgroup before exec, `pivot_root` after mounts).
- **Prerequisites:** P4–P7.
- **Observable result:** an Alpine shell inside limpet where Ctrl-C, Ctrl-Z, `fg`, and window resizing all behave. After `exit`, a check script finds no leftover mounts, cgroups, or processes.
- **Done when:** that holds and the threat model honestly lists what v0.1 doesn't protect against. Tag `v0.1`.
- **Extensions:** `limpet exec`, `limpet ps`, detached mode with logs, a state directory.
- **Leads to:** Part III, and P14's first prototype.

## Part III: Toward a real runtime

### P9 · Networking · L · Linux VM (root)
- **Build:** loopback-only networking by default, plus opt-in connectivity through a veth pair, a bridge, and NAT. Optionally, rootless networking through a user-mode network stack (`pasta` or `slirp4netns`).
- **Why:** For a sandbox running untrusted tools, the network is the main exfiltration path, and you can't control what you don't understand.
- **Learn.** *Go:* netlink (send one message by hand to bring `lo` up, then decide whether a library is worth it), the `net` package, a small HTTP CONNECT proxy for egress control. *Systems:* network namespaces, interfaces and routes, veth, bridges, NAT with nftables, DNS inside a root filesystem, egress filtering.
- **Prerequisites:** P8.
- **Observable result:** the sandbox reaches the internet only when enabled and only through what you allow, and teardown restores the VM's network exactly.
- **Done when:** you can trace a packet from a process in the sandbox to the internet and back, naming every hop and rule.
- **Extensions:** port forwarding, a domain allowlist in the proxy, per-sandbox bandwidth limits.
- **Leads to:** P14, where egress policy is central.

### P10 · Hardening: capabilities, seccomp, least privilege · L · Linux VM (root)
- **Build:** drop capabilities, set `no_new_privs`, install a seccomp-BPF filter assembled in Go, mask sensitive `/proc` and `/sys` paths, make the root read-only, and run as a non-root user inside. Then attack your own sandbox with known escape techniques and record what happens.
- **Why:** Namespaces change what a process sees; hardening limits what it can *do*. Many real container escapes exploit missing hardening rather than broken namespaces.
- **Learn.** *Go:* `encoding/binary`, building BPF programs, careful use of `unsafe` for syscalls that take struct pointers. *Systems:* the capability sets (permitted, effective, inheritable, bounding, ambient), setuid binaries and `no_new_privs`, seccomp filter semantics (including architecture checks: syscall numbers differ between arm64 and x86-64), how Docker's default seccomp profile works, an overview of LSMs (AppArmor, SELinux, Landlock).
- **Prerequisites:** P8.
- **Observable result:** probe shows the reduced capability sets and seccomp filter mode, and your escape attempts fail for the reasons you predicted.
- **Done when:** that holds and the threat model is updated, including what still isn't covered.
- **Extensions:** Landlock filesystem rules, seccomp user notification (brokering syscalls in a supervisor), rootless end to end.
- **Leads to:** P14.

### P11 · Speaking OCI · L · Linux VM
- **Build:** limpet accepts an OCI runtime bundle (`config.json` plus a root filesystem) and implements the runtime lifecycle (`create`, `start`, `state`, `kill`, `delete`) with a state directory. Optionally, it unpacks an OCI image into a bundle.
- **Why:** Standards are how real systems interoperate. Implementing a spec teaches precise reading, and shows the gap between "works for me" and "conforms".
- **Learn.** *Go:* JSON against a versioned schema, the runtime-spec Go types (your first large dependency), file locking, explicit state machines, `archive/tar`, `compress/gzip`, and `crypto/sha256` for content-addressed layers. *Systems:* the OCI runtime lifecycle and hooks, image layers and whiteouts, how containerd, Docker, and Podman delegate to runtimes such as runc.
- **Prerequisites:** P8 (P9–P10 help).
- **Observable result:** a bundle generated by a standard tool runs through `create` → `start` → `kill` → `delete` under limpet.
- **Done when:** that works, and you can list which parts of the spec limpet ignores and what breaks because of it.
- **Extensions:** run the runtime-tools validation suite, use limpet as Podman's `--runtime`, pull images from a registry.
- **Leads to:** P13c and P14 (one bundle format for every backend).

## Part IV: Virtual machines

Introduced only once Parts I–II are solid: a VM boundary is easiest to appreciate after building the container boundary by hand.

### P12 · Virtual machines from the outside · M–L · Linux VM with nested KVM
- **Build:** boot a Linux kernel in QEMU with an initramfs you build yourself, containing tinybox and `limpet init` as PID 1, so your code is the first process of a real VM. Talk to it over a serial console and virtio-vsock, and compare its isolation boundary with P8's container.
- **Why:** Hardware virtualization is the boundary that microVM sandboxes rely on. Before orchestrating or writing a VMM, you need to know what a guest needs in order to boot.
- **Learn.** *Go:* writing a `newc` cpio archive (the initramfs format), static cross-builds for the guest, `AF_VSOCK` sockets via `x/sys/unix`. *Systems:* what a hypervisor does, KVM's role vs. QEMU's, the boot flow (kernel, command line, initramfs, PID 1), virtio devices, the serial console, nested virtualization on Apple Silicon. Your M5 Pro is new enough for KVM inside the lab VM; we'll confirm it before relying on it.
- **Prerequisites:** P4 and P6 (P8 helps).
- **Observable result:** your initramfs boots to your init, which runs probe and a shell on the console and then powers the VM off cleanly.
- **Done when:** you can draw both boundaries (host → container process; host → VMM → guest kernel → guest process) and argue which attacks each one stops.
- **Extensions:** a minimal kernel configuration, measuring and cutting boot time, file sharing with virtio-fs.
- **Leads to:** P13.

### P13 · Choose a VMM path · L–XL · decision point
These are three distinct engineering goals. We'll choose one, or a sequence, when we get here, based on what you enjoyed and what P14 needs.
- **13a: Orchestrate an existing VMM.** Drive Firecracker (or Cloud Hypervisor) from Go through its HTTP API on a Unix socket: boot microVMs, attach root drives, run a guest agent over vsock, take snapshots. *Learn:* API clients over Unix sockets, supervising VMM processes, the jailer, guest/host protocols.
- **13b: Build a tiny KVM-based VMM.** Open `/dev/kvm`, create a VM, map guest memory, run a vCPU, and handle exits: first a few instructions, then serial output, then (stretch) a Linux boot. *Learn:* ioctls, `mmap`, `unsafe`, CPU registers, MMIO, device emulation. Most tutorials target x86-64; on Apple Silicon your guest is arm64, which changes vCPU setup, registers, and devices (there's no port I/O).
- **13c: One runtime, two backends.** Put containers and microVMs behind one limpet interface, designed around a real workload. *Learn:* interface design at a genuine seam, capability differences between backends, testing across them.
- **Prerequisites:** P12. **Done when:** defined at the decision point.

## Part V: The destination

### P14 · limpet sandbox: a home for dev tools and coding agents · XL
- **Build:** a sandbox for untrusted development tooling (a coding agent, `npm install`, a test suite). The project is mounted copy-on-write so you review a diff before anything touches the real files; the network is egress-filtered; there are resource and time limits; and an audit log records what ran.
- **Why:** This is the destination. It applies everything and forces an honest threat model: whose mistakes or malice are we containing, and with what confidence?
- **Learn:** threat modeling, policy design, overlay diffs, proxies, API design (a CLI plus a local socket), and how others solved the same problem: gVisor (a user-space kernel written in Go), bubblewrap, Firecracker-based sandboxes, Docker Sandboxes, and Anthropic's sandbox-runtime.
- **Prerequisites:** P8–P10, plus P13 if you choose a microVM backend.
- **Done when:** defined when we get there, from your threat model.

---

## Decision points

| When | Decision | Decided by |
|---|---|---|
| Start of P2 | How the Linux lab runs (default candidate: Lima; alternatives: OrbStack, UTM, a remote machine) | your preferences; nested-virtualization support for P12 |
| Start of P3 | Where limpet's code lives and how milestones are tagged | how P1–P2 code turned out |
| P5 | Rootful first, or rootless (user namespaces) from the start | how the first namespace experiments go |
| P9 | Hand-rolled netlink or a library | how the hand-built message goes |
| P13 | VMM path: 13a, 13b, 13c, or a sequence | interest, and what P14 needs |
| End of every project | Continue, detour, or reorder | the journal and `LEARNER.md` |

## Detours

Take any of these when curiosity strikes, and record it in `PROGRESS.md`.

- **D1 · tsh, a tiny shell** (after P3): pipelines, redirection, builtins (why must `cd` be one?), job control. The best preparation for P8's terminal work.
- **D2 · strace-lite** (after P3): trace a child's syscalls with `ptrace`.
- **D3 · The Go runtime underneath** (before P5): goroutine scheduling, OS threads, `LockOSThread`, and why multi-threaded processes and some namespace operations don't mix.
- **D4 · cpio and initramfs** (before P12): write the archive format yourself.
- **D5 · Landlock and macOS Seatbelt** (before P14): unprivileged sandboxing on both of your platforms.
- **D6 · Performance** (any time): `pprof`, `hyperfine`, `perf`; making `cat` and `yes` fast.
- **D7 · eBPF** (after P10): watch your sandbox from the kernel's side.

## Go topics: when each becomes necessary

| Topic | First needed |
|---|---|
| modules, packages, slices, strings, `switch`, errors, `defer` | P1 |
| interfaces (`io.Reader`/`io.Writer`), testing, `flag`, `encoding/json` | P1 |
| parsing, error wrapping, build constraints, cross-compilation | P2 |
| goroutines, channels, `select`, `os/signal`, `context` | P3 |
| concurrency patterns, `sync`, the race detector, `log/slog` | P4 |
| `SysProcAttr`, re-exec, `runtime.LockOSThread`, `x/sys/unix` | P5 |
| package design, `internal/`, end-to-end tests | P8 |
| `encoding/binary`, `unsafe` (carefully) | P10 |
| third-party dependencies, `archive/tar`, hashing | P11 |
| ioctls, `mmap`, raw memory (if 13b) | P13 |
| generics | when a real duplication problem calls for them, not before |

## What limpet is not

- **Not a secure sandbox,** at any stage, until a threat model and evidence say otherwise for a specific, stated threat.
- **Not a replacement** for Docker, runc, or Firecracker. Those are reference points for reading after you've built your own version.
- **Not the point.** The point is the person who can build it.
