# Experiments

Disposable investigations: a quick program to test a hunch, a reproduction of something surprising, a benchmark. Use one folder per experiment, named `YYYY-MM-DD-topic/`.

Each folder gets a short `NOTES.md`: **question → prediction → what happened → explanation**. The prediction is the part that teaches, so write it before running anything.

A single-file Go program needs no module: `go run main.go` works as long as it imports only the standard library.

Privileged experiments (mounts, namespaces, cgroups) run only inside the Linux lab VM, never on macOS, and their notes must include cleanup steps.
