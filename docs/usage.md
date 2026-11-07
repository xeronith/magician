# Usage

## `magician version`

Prints the version (default `dev` when built locally; injected by the release pipeline).

```
$ magician version
dev
```

## `magician greet [name]`

Prints a friendly greeting. Without a name (or with an empty/whitespace-only name) it greets `world`.

```
$ magician greet
Hello, world!

$ magician greet Alice
Hello, Alice!

$ magician greet "  Bob  "
Hello, Bob!
```

## `magician spell <word>`

Prints the word reversed. Exits with a non-zero status if `<word>` is missing.

```
$ magician spell hello
olleh

$ magician spell
Error: spell requires a <word> argument
```
