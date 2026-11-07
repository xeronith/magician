# Task: review this pull request

You are reviewing an automated change to a Go repository. The unified diff is attached. Read files in the repository when you need more context.

## Blocking issues (VERDICT: REQUEST_CHANGES)

- The change does not implement what `PROMPT.md` asks for.
- A correctness bug, panic, data race, or unhandled error on a normal path.
- New behavior has no meaningful test, or an existing test was weakened.
- A security issue, leaked secret, or a change to a protected path.
- The repository does not build or the tests are clearly broken.

## Non-blocking

Naming, formatting, comments, minor style, and extra polish are not reasons to block a merge.

## Output

Write your review to `REVIEW.md` at the repository root. The first line must be exactly:

```text
VERDICT: APPROVE
```

or:

```text
VERDICT: REQUEST_CHANGES
```

Follow the verdict with a short markdown review. For blocking findings, include the file and line, for example:

```md
- [high] cmd/magician/main.go:42 - the empty input path panics; return an error instead.
```

Keep the review under 300 words. Do not modify any file other than `REVIEW.md`. Do not run shell commands.
