package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGreeting(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty name becomes world", input: "", want: "Hello, world!"},
		{name: "whitespace-only name becomes world", input: "   \t\n", want: "Hello, world!"},
		{name: "surrounding whitespace is trimmed", input: "\t Alice \n", want: "Hello, Alice!"},
		{name: "plain name is kept", input: "Alice", want: "Hello, Alice!"},
		{name: "inner whitespace is kept", input: "  Alice Smith  ", want: "Hello, Alice Smith!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := greeting(tt.input); got != tt.want {
				t.Errorf("greeting(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty word", input: "", want: ""},
		{name: "single rune", input: "a", want: "a"},
		{name: "ascii word", input: "magician", want: "naicigam"},
		{name: "palindrome", input: "racecar", want: "racecar"},
		{name: "utf-8 runes", input: "héllo", want: "olléh"},
		{name: "emoji counts as one rune", input: "Go🚀", want: "🚀oG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.input); got != tt.want {
				t.Errorf("reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestShout(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "single word is uppercased", input: "magic", want: "MAGIC!!!"},
		{name: "multiple words are uppercased", input: "hello there", want: "HELLO THERE!!!"},
		{name: "mixed case becomes upper", input: "AbC dEf", want: "ABC DEF!!!"},
		{name: "empty text still shouts", input: "", want: "!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shout(tt.input); got != tt.want {
				t.Errorf("shout(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "version prints the version", args: []string{"version"}, want: "dev\n"},
		{name: "greet without a name greets world", args: []string{"greet"}, want: "Hello, world!\n"},
		{name: "greet with a name", args: []string{"greet", "Alice"}, want: "Hello, Alice!\n"},
		{name: "spell reverses the word", args: []string{"spell", "magic"}, want: "cigam\n"},
		{name: "no command is a usage error", args: nil, wantErr: true},
		{name: "unknown command is an error", args: []string{"fly"}, wantErr: true},
		{name: "spell without a word is an error", args: []string{"spell"}, wantErr: true},
		{name: "shout joins and uppercases words", args: []string{"shout", "hello", "world"}, want: "HELLO WORLD!!!\n"},
		{name: "shout without text is an error", args: []string{"shout"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(&out, tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("run(%v) succeeded, want error", tt.args)
				}
				if strings.TrimSpace(out.String()) != "" {
					t.Errorf("run(%v) wrote %q alongside an error", tt.args, out.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("run(%v) returned error %v", tt.args, err)
			}
			if got := out.String(); got != tt.want {
				t.Errorf("run(%v) wrote %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
