# Usage

`magician` is a small, dependency-free CLI. Build it locally with:

```console
$ go build -o magician ./cmd/magician
```

## magician version

Prints the version. Binaries built without the release pipeline report `dev`; releases inject the tag with `-ldflags "-X main.version=<tag>"`.

```console
$ magician version
dev
```

## magician greet [name]

Prints a friendly greeting. The name is optional: when it is omitted, empty, or only whitespace, the greeting uses `world`. Surrounding whitespace is trimmed.

```console
$ magician greet
Hello, world!

$ magician greet Alice
Hello, Alice!

$ magician greet "  Alice  "
Hello, Alice!
```

## magician spell &lt;word&gt;

Prints the word reversed, by Unicode code points.

```console
$ magician spell magic
cigam

$ magician spell héllo
olléh
```

When the word is missing, the command exits with a non-zero status and a clear error:

```console
$ magician spell
spell: missing word

$ echo $?
1
```

## magician shout &lt;text&gt;

Prints the text uppercased, followed by `!!!`. Multiple arguments are joined with a single space.

```console
$ magician shout magic
MAGIC!!!

$ magician shout hello world
HELLO WORLD!!!
```

When the text is missing, the command exits with a non-zero status and a clear error:

```console
$ magician shout
shout: missing text

$ echo $?
1
```
