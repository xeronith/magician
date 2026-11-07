// Command magician is a small CLI with version, greet, spell, shout, and count commands.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

var version = "dev"

// errUsage describes the accepted commands.
var errUsage = errors.New("usage: magician <version|greet|spell|shout|count> [args]")

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run dispatches args to the matching command and writes output to out.
func run(out io.Writer, args []string) error {
	if len(args) == 0 {
		return errUsage
	}

	switch args[0] {
	case "version":
		say(out, version)
		return nil
	case "greet":
		say(out, greeting(strings.Join(args[1:], " ")))
		return nil
	case "spell":
		if len(args) < 2 {
			return errors.New("spell: missing word")
		}
		say(out, reverse(args[1]))
		return nil
	case "shout":
		if len(args) < 2 {
			return errors.New("shout: missing text")
		}
		say(out, shout(strings.Join(args[1:], " ")))
		return nil
	case "count":
		if len(args) < 2 {
			return errors.New("count: missing text")
		}
		say(out, fmt.Sprintf("%d", count(strings.Join(args[1:], " "))))
		return nil
	default:
		return fmt.Errorf("unknown command %q: %w", args[0], errUsage)
	}
}

// say writes one line to out. The write error is discarded because there
// is nothing left to do once the command has produced its output.
func say(out io.Writer, text string) {
	_, _ = fmt.Fprintln(out, text)
}

// greeting returns a friendly greeting for name, using "world" when the
// name is empty or contains only whitespace.
func greeting(name string) string {
	if name = strings.TrimSpace(name); name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// reverse returns s with its runes in reverse order.
func reverse(s string) string {
	runes := []rune(s)
	slices.Reverse(runes)
	return string(runes)
}

// shout returns text uppercased and followed by "!!!".
func shout(text string) string {
	return strings.ToUpper(text) + "!!!"
}

// count returns the number of whitespace-separated words in text.
func count(text string) int {
	return len(strings.Fields(text))
}
