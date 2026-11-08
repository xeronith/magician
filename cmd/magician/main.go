// Command magician is a small CLI with version, greet, spell, shout, count,
// sort, and palindrome commands, plus serve: an HTTP server with a small spell
// REST API and Scalar documentation.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
)

var version = "dev"

// errUsage describes the accepted commands.
var errUsage = errors.New("usage: magician <version|greet|spell|shout|count|sort|palindrome|serve> [args]")

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
	case "sort":
		if len(args) < 2 {
			return errors.New("sort: missing text")
		}
		say(out, strings.Join(sorted(strings.Join(args[1:], " ")), " "))
		return nil
	case "palindrome":
		if len(args) < 2 {
			return errors.New("palindrome: missing text")
		}
		say(out, strconv.FormatBool(palindrome(strings.Join(args[1:], " "))))
		return nil
	case "serve":
		return serve(out, args[1:])
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

// sorted returns the whitespace-separated words of text in ascending byte
// order, keeping duplicates.
func sorted(text string) []string {
	words := strings.Fields(text)
	slices.Sort(words)
	return words
}

// palindrome reports whether text reads the same forwards and backwards
// ignoring case, whitespace, and punctuation: only letters and digits are
// compared. Text with no letters or digits counts as a palindrome.
func palindrome(text string) bool {
	var cleaned []rune
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
		if cleaned[i] != cleaned[j] {
			return false
		}
	}
	return true
}

// defaultAddr is the address serve listens on when none is given.
const defaultAddr = "127.0.0.1:8080"

// serve takes the optional address and runs the HTTP server until the process
// is interrupted.
func serve(out io.Writer, args []string) error {
	addr := defaultAddr
	switch len(args) {
	case 0:
	case 1:
		addr = args[0]
	default:
		return fmt.Errorf("serve: unexpected argument %q", args[1])
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return serveAddr(out, addr, ctx)
}

// serveAddr listens on addr and serves requests until ctx is interrupted or
// the listener fails.
func serveAddr(out io.Writer, addr string, ctx context.Context) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	server := &http.Server{Handler: newSpellHandler()}
	errCh := make(chan error, 1)
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	say(out, "listening on http://"+listener.Addr().String())

	select {
	case err := <-errCh:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	// The process was interrupted; give in-flight requests a moment to
	// finish before returning.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	return nil
}
