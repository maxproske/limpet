# Resources

Curated, not exhaustive. Each entry says what the resource is and when it's useful. Authoritative sources (Go docs, man pages, kernel docs, specs) outrank tutorials, and tutorials and example code are reference material to read critically.

**Spoilers.** Entries tagged **[spoiler: after Pn]** contain complete solutions to course exercises. Read them after you've attempted project Pn, or whenever you deliberately choose to.

## Inspiration

- **Handmade Hero**, Casey Muratori. https://handmadehero.org/ now redirects to Molly Rocket and says the series has been retired; the [episode guide](https://guide.handmadehero.org/) is still up. A complete game built from scratch in C on Windows, live, one day at a time. It inspires the *method* of this course: build everything, understand everything, keep a daily history. Casey's current material on performance-aware programming is at [Computer, Enhance!](https://www.computerenhance.com/).
- **Games With Go**, Jack Mott. https://gameswithgo.org/ is a beginner video series that teaches programming through small Go games (SDL2, Pong, noise, a genetic-art evolver). The site looks unmaintained: its outbound links are gone and it carries a hidden spam link. The companion code is at [jackmott/gameswithgo-public](https://github.com/jackmott/gameswithgo-public) (2019, from before Go modules), and your own run-through from 2018 is at [maxproske/games-with-go](https://github.com/maxproske/games-with-go).

## Go

**Official**

- [Go documentation](https://go.dev/doc/): start here when you're unsure where something lives.
- [A Tour of Go](https://go.dev/tour/): a fast syntax refresher.
- [How to Write Go Code](https://go.dev/doc/code): modules, packages, and the `go` command (P1).
- [Organizing a Go module](https://go.dev/doc/modules/layout): for when tinybox grows packages, and for P3's layout decision.
- [Managing dependencies](https://go.dev/doc/modules/managing-dependencies): before your first dependency (P2/P5).
- [The Go Programming Language Specification](https://go.dev/ref/spec): the final word on what code means.
- [The Go Memory Model](https://go.dev/ref/mem): P4, once goroutines share data.
- [Effective Go](https://go.dev/doc/effective_go): core idioms. It dates from 2009 and has barely changed, so it doesn't cover modules, generics, or newer packages.
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments): the checklist Go reviewers use. Skim it before your first review.
- [Google's Go Style Guide](https://google.github.io/styleguide/go/): style with rationale, when you want the "why".
- [Release history](https://go.dev/doc/devel/release): what changed since you last wrote Go regularly.
- Go blog: [Errors are values](https://go.dev/blog/errors-are-values) (P1 L2), [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) (P2), [Pipelines and cancellation](https://go.dev/blog/pipelines) (P4).
- [Data Race Detector](https://go.dev/doc/articles/race_detector): P3 and P4.

**Packages you'll live in.** `go doc` shows the same text in your terminal.
[`os`](https://pkg.go.dev/os) · [`os/exec`](https://pkg.go.dev/os/exec) · [`os/signal`](https://pkg.go.dev/os/signal) · [`syscall`](https://pkg.go.dev/syscall) · [`golang.org/x/sys/unix`](https://pkg.go.dev/golang.org/x/sys/unix) · [`runtime.LockOSThread`](https://pkg.go.dev/runtime#LockOSThread) (P5) · the full [standard library](https://pkg.go.dev/std)

**Books and references**

- [The Go Programming Language](https://www.gopl.io/), Donovan and Kernighan. The classic. It predates modules and generics, but its fundamentals and exercises are excellent (you started it in 2018).
- [Learning Go, 2nd edition](https://www.oreilly.com/library/view/learning-go-2nd/9781098139285/), Jon Bodner. Modern Go (you started it this year). A good companion for looking up language features as they come up.
- [Go by Example](https://gobyexample.com/): annotated snippets for quick "how do I…" questions.
- [100 Go Mistakes and How to Avoid Them](https://100go.co/): most useful after P3, once you have real code to compare.

## Unix and Linux fundamentals

- **Linux man pages** at [man7.org](https://man7.org/linux/man-pages/index.html), the authoritative reference. Section 2 covers system calls and section 7 covers overviews. `man` on macOS gives BSD versions, so for Linux behavior use man7.org or `man` inside the lab VM. Key pages by project:
  - P1–P3: [execve(2)](https://man7.org/linux/man-pages/man2/execve.2.html), [fork(2)](https://man7.org/linux/man-pages/man2/fork.2.html), [wait(2)](https://man7.org/linux/man-pages/man2/wait.2.html), [open(2)](https://man7.org/linux/man-pages/man2/open.2.html), [read(2)](https://man7.org/linux/man-pages/man2/read.2.html), [pipe(7)](https://man7.org/linux/man-pages/man7/pipe.7.html), [signal(7)](https://man7.org/linux/man-pages/man7/signal.7.html)
  - P2: [proc(5)](https://man7.org/linux/man-pages/man5/proc.5.html), [syscalls(2)](https://man7.org/linux/man-pages/man2/syscalls.2.html), [strace(1)](https://man7.org/linux/man-pages/man1/strace.1.html)
  - P3–P4: [prctl(2)](https://man7.org/linux/man-pages/man2/prctl.2.html)
- [The Linux Programming Interface](https://man7.org/tlpi/), Michael Kerrisk, who also maintains the man pages. *The* systems-programming book for this course; its chapters map onto P2–P10.
- [Advanced Programming in the UNIX Environment](http://www.apuebook.com/), Stevens and Rago. Classic Unix rather than Linux-specific: processes, signals, terminals. (The site is HTTP-only.)
- [Operating Systems: Three Easy Pieces](https://pages.cs.wisc.edu/~remzi/OSTEP/), a free textbook. Read the process-API chapters alongside P1–P3, concurrency before P4, and virtualization before P12.
- [Julia Evans](https://jvns.ca/) and [Wizard Zines](https://wizardzines.com/): friendly explanations of strace, `/proc`, and containers. Good warm-ups before a new topic.
- [The TTY demystified](https://www.linusakesson.net/programming/tty/), Linus Åkesson. Terminals, sessions, and job control; essential for detour D1 and for P8.
- [Implementing a Shell](https://sourceware.org/glibc/manual/latest/html_node/Implementing-a-Shell.html), from the GNU C Library manual. Job control done right; for D1 and P8.

## Containers

**Authoritative**

- Namespaces: [namespaces(7)](https://man7.org/linux/man-pages/man7/namespaces.7.html), [pid_namespaces(7)](https://man7.org/linux/man-pages/man7/pid_namespaces.7.html), [mount_namespaces(7)](https://man7.org/linux/man-pages/man7/mount_namespaces.7.html), [user_namespaces(7)](https://man7.org/linux/man-pages/man7/user_namespaces.7.html), [network_namespaces(7)](https://man7.org/linux/man-pages/man7/network_namespaces.7.html), [cgroup_namespaces(7)](https://man7.org/linux/man-pages/man7/cgroup_namespaces.7.html), [clone(2)](https://man7.org/linux/man-pages/man2/clone.2.html), [unshare(2)](https://man7.org/linux/man-pages/man2/unshare.2.html), [setns(2)](https://man7.org/linux/man-pages/man2/setns.2.html). For P5.
- Filesystems: [mount(2)](https://man7.org/linux/man-pages/man2/mount.2.html), [pivot_root(2)](https://man7.org/linux/man-pages/man2/pivot_root.2.html), [chroot(2)](https://man7.org/linux/man-pages/man2/chroot.2.html), and the kernel's [overlayfs docs](https://docs.kernel.org/filesystems/overlayfs.html). For P6.
- [Control Group v2](https://docs.kernel.org/admin-guide/cgroup-v2.html) (kernel docs) and [cgroups(7)](https://man7.org/linux/man-pages/man7/cgroups.7.html). The authoritative cgroup reference for P7. Avoid old version-pinned copies of the kernel docs, which lack newer files such as `cgroup.kill`.
- Hardening: [capabilities(7)](https://man7.org/linux/man-pages/man7/capabilities.7.html), [seccomp(2)](https://man7.org/linux/man-pages/man2/seccomp.2.html), the kernel's [Seccomp BPF](https://docs.kernel.org/userspace-api/seccomp_filter.html) and [Landlock](https://docs.kernel.org/userspace-api/landlock.html) docs. For P10 and D5.
- Plumbing: [veth(4)](https://man7.org/linux/man-pages/man4/veth.4.html) for P9, [pty(7)](https://man7.org/linux/man-pages/man7/pty.7.html) for P8.
- [OCI runtime specification](https://github.com/opencontainers/runtime-spec), latest release v1.3.0 (November 2025). The contract between Docker, containerd, or Podman and runtimes such as runc. Read `bundle.md`, `runtime.md`, `runtime-linux.md`, `config.md`, and `config-linux.md` first. Skim it after P8 to see how much you've already built, then implement it in P11.
- [OCI image specification](https://github.com/opencontainers/image-spec): layers, manifests, and digests. A P11 extension.

**Explanations**

- [What even is a container?](https://jvns.ca/blog/2016/10/10/what-even-is-a-container/), Julia Evans. Short; worth reading now.
- [Namespaces in operation](https://lwn.net/Articles/531114/), Michael Kerrisk on LWN. A 2013 series that's still the best conceptual introduction to namespaces. Read it at P5, and check details against the current man pages.
- [Learning Containers From The Bottom Up](https://iximiuz.com/en/posts/container-learning-path/), Ivan Velichko. A curated path that runs parallel to this course. See also his [blog archive](https://iximiuz.com/en/archive/) and [hands-on labs](https://labs.iximiuz.com/).
- [Container Security](https://www.oreilly.com/library/view/container-security/9781492056690/), Liz Rice. Organized around threat models, for P10 and P14. That link is the first edition; a second edition came out in 2025.

**Reference implementations**

- [Containers From Scratch](https://www.youtube.com/watch?v=8fi7uSYlOdc), Liz Rice (GOTO 2018), plus [the code](https://github.com/lizrice/containers-from-scratch). A 40-minute live-coded container in Go. **[spoiler: after P6]**, though watching only the first ten minutes for motivation is fine.
- [Linux containers in 500 lines of code](https://blog.lizzie.io/linux-containers-in-500-loc.html), Lizzie Dixon. Written in C and thorough about security details. **[spoiler: after P10]**
- [runc](https://github.com/opencontainers/runc), the reference OCI runtime, written in Go. After P8, read specific parts to compare with your design rather than reading it front to back. **[spoiler: after P8]**
- [tini](https://github.com/krallin/tini), the tiny init behind `docker run --init`. **[spoiler: after P4]**
- [runtime-tools](https://github.com/opencontainers/runtime-tools), the OCI validation suite, for P11.

**The tutorial series you found.** Read these critically. They're useful, but they contain real errors.

- [Reddit discussion](https://www.reddit.com/r/golang/comments/1spj2bu/), April 2026: "I built a container from scratch in Go to understand what Docker actually does under the hood", posted by the author of the posts below. The comments are worth more than the post: they point to Liz Rice's talk, chroot escapes and `pivot_root`, and the OCI runtime spec. (Reddit blocks automated clients, but the thread loads in a browser.)
- [Containers Are Not Magic, Part 1: namespaces](https://blog.iamvedant.in/containers-are-not-magic-namespaces-from-scratch). About 100 lines of Go: re-exec through `/proc/self/exe` into new UTS, PID, and mount namespaces, private mount propagation, and a fresh `/proc`. The core mechanics are right, but several explanations aren't. **[spoiler: after P5]**
- [Containers Are Not Magic, Part 2: cgroups and chroot](https://blog.iamvedant.in/containers-are-not-magic-cgroups-and-chroot-from-scratch). cgroup v2 memory and pid limits, `chroot` into Alpine, and an OOM demo. It contains significant errors about chroot's security, cleanup, and when a process lands in its cgroup. **[spoiler: after P7]**
- [Vedant-Gandhi/How-containers-work](https://github.com/Vedant-Gandhi/How-containers-work), the posts' code. A complete, and buggy, version of P5–P7. Its setup script downloads an x86-64 Alpine 3.19 root filesystem, which is end-of-life and won't run in your arm64 VM. **[spoiler: after P7]**
- After P7 you'll review this series as if it were a pull request. The errata are kept in [`spoilers/containers-are-not-magic.md`](spoilers/containers-are-not-magic.md) for comparing afterwards, so leave that file closed until then.

## Virtual machines

- [Using the KVM API](https://lwn.net/Articles/658511/), Josh Triplett on LWN (2015). The smallest complete KVM program. It's x86-specific; an arm64 guest on your Mac needs a different vCPU setup.
- [The KVM API](https://docs.kernel.org/virt/kvm/api.html) (kernel docs). Authoritative and dense; path 13b.
- [QEMU documentation](https://www.qemu.org/docs/master/), for P12.
- [Firecracker](https://firecracker-microvm.github.io/): [repository](https://github.com/firecracker-microvm/firecracker), [design document](https://github.com/firecracker-microvm/firecracker/blob/main/docs/design.md), and the NSDI 2020 paper, [Firecracker: Lightweight Virtualization for Serverless Applications](https://www.usenix.org/conference/nsdi20/presentation/agache). For P13. Its [Go SDK](https://github.com/firecracker-microvm/firecracker-go-sdk) is **[spoiler: after 13a]**.
- [Cloud Hypervisor](https://www.cloudhypervisor.org/), an alternative VMM for path 13a.
- [rust-vmm](https://github.com/rust-vmm): shared building blocks for VMMs, and good reading for 13b.
- [Kata Containers](https://katacontainers.io/): containers inside lightweight VMs, which is idea 13c in production.
- [apple/containerization](https://github.com/apple/containerization): runs each Linux container in its own lightweight VM, using Apple's Virtualization framework. A design worth comparing in 13c and P14.
- [Lima](https://lima-vm.io/), your lab-VM manager from P2 on. Nested virtualization, which P12 needs for KVM inside the lab, has been supported since Lima v1.0. It requires an M3 or later chip and macOS 15 or later, enabled with `nestedVirtualization: true` or `limactl create --nested-virt`.

## Sandboxing (P14)

- [gVisor](https://gvisor.dev/) ([source](https://github.com/google/gvisor)): an application kernel written in Go. After P10, it's a master class in intercepting system calls.
- [bubblewrap](https://github.com/containers/bubblewrap): unprivileged sandboxing with namespaces, used by Flatpak and by sandbox-runtime.
- [Docker Sandboxes](https://docs.docker.com/ai/sandboxes/): Docker's sandbox for coding agents, which runs each agent in its own microVM with its own kernel. The production answer to P14's question.
- [sandbox-runtime](https://github.com/anthropics/sandbox-runtime), from Anthropic: bubblewrap on Linux and Seatbelt on macOS, with network access forced through proxies that enforce domain allowlists. A design without containers to compare with yours.

---

_Link check, 2026-10-13: every link above loaded except O'Reilly's pages and Reddit, which block automated clients and load normally in a browser._
