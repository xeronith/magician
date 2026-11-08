package main

import (
	"bytes"
	"context"
	"slices"
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

func TestCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "empty text has no words", input: "", want: 0},
		{name: "whitespace-only text has no words", input: "  \t\n ", want: 0},
		{name: "single word", input: "magic", want: 1},
		{name: "multiple words", input: "hello big world", want: 3},
		{name: "runs of whitespace separate words", input: "  one\ttwo\nthree  ", want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := count(tt.input); got != tt.want {
				t.Errorf("count(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty text has no words", input: "", want: nil},
		{name: "whitespace-only text has no words", input: "  \t\n ", want: nil},
		{name: "a single word", input: "magic", want: []string{"magic"}},
		{name: "multiple words are sorted", input: "banana apple cherry", want: []string{"apple", "banana", "cherry"}},
		{name: "runs of whitespace separate words", input: "  one\ttwo\nthree  ", want: []string{"one", "three", "two"}},
		{name: "uppercase sorts before lowercase", input: "zeta Alpha beta", want: []string{"Alpha", "beta", "zeta"}},
		{name: "duplicates are kept", input: "b a b a", want: []string{"a", "a", "b", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sorted(tt.input); !slices.Equal(got, tt.want) {
				t.Errorf("sorted(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "empty text is a palindrome", input: "", want: true},
		{name: "whitespace-only text is a palindrome", input: "  \t\n ", want: true},
		{name: "punctuation-only text is a palindrome", input: "!?", want: true},
		{name: "a palindrome word", input: "racecar", want: true},
		{name: "case is ignored", input: "RaceCar", want: true},
		{name: "whitespace and punctuation are ignored", input: "A man, a plan, a canal: Panama!", want: true},
		{name: "not a palindrome", input: "hello", want: false},
		{name: "digits are compared", input: "12321", want: true},
		{name: "digits that differ", input: "12345", want: false},
		{name: "unicode letters are compared", input: "hél léh", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := palindrome(tt.input); got != tt.want {
				t.Errorf("palindrome(%q) = %t, want %t", tt.input, got, tt.want)
			}
		})
	}
}

func TestServe(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "an unusable address is an error",
			args: []string{"serve", "bogus"},
			want: "serve: listen tcp: address bogus: missing port in address",
		},
		{
			name: "a second argument is an error",
			args: []string{"serve", "127.0.0.1:8080", "extra"},
			want: `serve: unexpected argument "extra"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(&out, tt.args)
			if err == nil {
				t.Fatalf("run(%v) succeeded, want error", tt.args)
			}
			if err.Error() != tt.want {
				t.Errorf("run(%v) returned %q, want %q", tt.args, err, tt.want)
			}
			if out.Len() != 0 {
				t.Errorf("run(%v) wrote %q alongside an error", tt.args, out.String())
			}
		})
	}
}

func TestServeAddrShutsDownOnInterrupt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var out bytes.Buffer
	if err := serveAddr(&out, "127.0.0.1:0", ctx); err != nil {
		t.Fatalf("serveAddr returned error: %v", err)
	}
	if !strings.HasPrefix(out.String(), "listening on http://127.0.0.1:") {
		t.Errorf("serveAddr wrote %q, want a listening message", out.String())
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
		{name: "count counts the words", args: []string{"count", "hello", "big", "world"}, want: "3\n"},
		{name: "count without text is an error", args: []string{"count"}, wantErr: true},
		{name: "sort joins and sorts the words", args: []string{"sort", "banana", "apple", "cherry"}, want: "apple banana cherry\n"},
		{name: "sort sorts by byte order", args: []string{"sort", "zeta", "Alpha"}, want: "Alpha zeta\n"},
		{name: "sort keeps duplicates", args: []string{"sort", "b", "a", "b"}, want: "a b b\n"},
		{name: "sort of whitespace-only text prints an empty line", args: []string{"sort", "   "}, want: "\n"},
		{name: "sort without text is an error", args: []string{"sort"}, wantErr: true},
		{name: "palindrome joins and checks the words", args: []string{"palindrome", "Never", "odd", "or", "even"}, want: "true\n"},
		{name: "palindrome reports false", args: []string{"palindrome", "hello"}, want: "false\n"},
		{name: "palindrome without text is an error", args: []string{"palindrome"}, wantErr: true},
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
