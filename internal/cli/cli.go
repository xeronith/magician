// Package cli implements the commands of the magician CLI.
package cli

import (
	"fmt"
	"strings"
)

// Greet returns a friendly greeting for the given name.
// An empty or whitespace-only name is treated as "world".
func Greet(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Spell returns the word reversed.
func Spell(word string) string {
	runes := []rune(word)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
