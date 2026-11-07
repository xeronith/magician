package cli

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "Hello, world!"},
		{"whitespace only", "   \t ", "Hello, world!"},
		{"name", "Alice", "Hello, Alice!"},
		{"trim", "  Bob  ", "Hello, Bob!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Greet(tt.in); got != tt.want {
				t.Errorf("Greet(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSpell(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"simple", "abc", "cba"},
		{"unicode", "hällo", "olläh"},
		{"palindrome", "anna", "anna"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Spell(tt.in); got != tt.want {
				t.Errorf("Spell(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
