// Command magician is a small CLI tool.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/xeronith/magician/internal/cli"
)

var version = "dev"

func usage(w io.Writer) {
	fmt.Fprintf(w, "Usage: magician <command> [args]\n\nCommands:\n  version      print the version\n  greet [name] print a friendly greeting\n  spell <word> print the word reversed\n")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage(os.Stderr)
		os.Exit(2)
	}

	switch args[0] {
	case "version":
		fmt.Println(version)
	case "greet":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		fmt.Println(cli.Greet(name))
	case "spell":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: spell requires a <word> argument")
			os.Exit(1)
		}
		fmt.Println(cli.Spell(args[1]))
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", args[0])
		usage(os.Stderr)
		os.Exit(2)
	}
}
