# Usage

`magician` is a small, dependency-free CLI. Build it locally with:

```console
$ go build -o magician ./cmd/magician
```

Besides the commands below, `magician serve` runs an HTTP server with a small spell REST API and Scalar documentation.

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

## magician count &lt;text&gt;

Prints the number of whitespace-separated words in the text. Multiple arguments are joined with a single space. Whitespace-only text prints `0`.

```console
$ magician count hello world
2

$ magician count "  one   two  "
2

$ magician count "   "
0
```

When the text is missing, the command exits with a non-zero status and a clear error:

```console
$ magician count
count: missing text

$ echo $?
1
```

## magician sort &lt;text&gt;

Prints the whitespace-separated words sorted in ascending order and joined by a single space. Multiple arguments are joined with a single space. Sorting is by byte order (Go string comparison), so uppercase letters come before lowercase ones. Duplicates are kept. Text with no words prints an empty line.

```console
$ magician sort banana apple cherry
apple banana cherry

$ magician sort "one   two  three"
one three two

$ magician sort zeta Alpha beta
Alpha beta zeta

$ magician sort b a b a
a a b b

$ magician sort "   "

```

When the text is missing, the command exits with a non-zero status and a clear error:

```console
$ magician sort
sort: missing text

$ echo $?
1
```

## magician palindrome &lt;text&gt;

Prints `true` when the text reads the same forwards and backwards, otherwise `false`. Case, whitespace, and punctuation are ignored: only letters and digits are compared. Multiple arguments are joined with a single space. Text with no letters or digits, including whitespace-only text, prints `true`.

```console
$ magician palindrome racecar
true

$ magician palindrome "Never odd or even"
true

$ magician palindrome "A man, a plan, a canal: Panama!"
true

$ magician palindrome hello
false

$ magician palindrome "  12321  "
true

$ magician palindrome "   "
true
```

When the text is missing, the command exits with a non-zero status and a clear error:

```console
$ magician palindrome
palindrome: missing text

$ echo $?
1
```

## magician serve [addr]

Starts an HTTP server and blocks until the process is interrupted with `Ctrl+C`. The address is optional and defaults to `127.0.0.1:8080`.

```console
$ magician serve
listening on http://127.0.0.1:8080
```

A custom address can be passed:

```console
$ magician serve 127.0.0.1:9090
listening on http://127.0.0.1:9090
```

When the address cannot be used, the command exits with a non-zero status and a clear error:

```console
$ magician serve bogus
serve: listen tcp: address bogus: missing port in address

$ echo $?
1
```

## REST API

`magician serve` exposes a small spell catalog. Every response with a body is JSON. Spells live in memory and are reset when the process exits; ids start at 1, increase on every creation, and are never reused. The examples below assume the server runs on `127.0.0.1:8080`.

### GET /health

```console
$ curl http://127.0.0.1:8080/health
{"status":"ok"}
```

### GET /api/spells

Lists every spell, oldest first, as a JSON array:

```console
$ curl http://127.0.0.1:8080/api/spells
[]
```

### POST /api/spells

Creates a spell. Both fields are required; surrounding whitespace is trimmed, and a blank field is rejected with `400`.

```console
$ curl -X POST -d '{"name":"  fireball  ","effect":"burns"}' http://127.0.0.1:8080/api/spells
{"id":1,"name":"fireball","effect":"burns"}
```

Invalid JSON or a blank field returns `400`:

```console
$ curl -X POST -d '{"name":"","effect":"burns"}' http://127.0.0.1:8080/api/spells
{"error":"name must not be blank"}

$ curl -X POST -d 'not json' http://127.0.0.1:8080/api/spells
{"error":"invalid JSON"}
```

### GET /api/spells/{id}

```console
$ curl http://127.0.0.1:8080/api/spells/1
{"id":1,"name":"fireball","effect":"burns"}

$ curl http://127.0.0.1:8080/api/spells/999
{"error":"spell not found"}
```

### DELETE /api/spells/{id}

Deletes the spell and answers `204` with no body:

```console
$ curl -i -X DELETE http://127.0.0.1:8080/api/spells/1
HTTP/1.1 204 No Content

$ curl -i -X DELETE http://127.0.0.1:8080/api/spells/999
{"error":"spell not found"}
```

### GET /openapi.json

Serves the OpenAPI 3.1 document describing every route and schema:

```console
$ curl http://127.0.0.1:8080/openapi.json
{"openapi":"3.1.0","info":{"title":"magician",...}}
```

The document's `servers` entry is the relative URL `/`, so Scalar calls the same origin that served `/docs`. The reference UI works through `localhost`, `127.0.0.1`, or any custom address without CORS.

### GET /docs

Serves an HTML page that renders `/openapi.json` with Scalar. Open `http://127.0.0.1:8080/docs` in a browser. The page imports Scalar from the pinned URL `https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.77.0/esm.js`; the script runs in the browser only, and the server never fetches the CDN itself.

### GET /

Redirects to `/docs`:

```console
$ curl -i http://127.0.0.1:8080/
HTTP/1.1 302 Found
Location: /docs
```
