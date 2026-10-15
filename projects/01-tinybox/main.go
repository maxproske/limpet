// Tinybox is a single binary containing small Unix tools (applets).
//
// Usage:
//
//	tinybox <applet> [args...]
//
// Applets: true, false, echo [-n] [words...].
// Usage errors exit with status 2, the convention Go's flag package follows.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tinybox <applet> [args...]")
		os.Exit(2)
	}
	applet, args := os.Args[1], os.Args[2:]

	switch applet {
	case "true":
		// Could drop this since main exits with 0.
		os.Exit(0)
	case "false":
		os.Exit(1)
	case "echo":
		if len(args) > 0 && args[0] == "-n" {
			fmt.Print(strings.Join(args[1:], " "))
		} else {
			fmt.Println(strings.Join(args, " "))
		}
	default:
		fmt.Fprintf(os.Stderr, "tinybox: unknown applet %q\n", applet)
		os.Exit(2)
	}
}
