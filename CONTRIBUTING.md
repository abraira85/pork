# Contributing to Pork

Thanks for wanting to make Pork better. Every contribution counts — code, docs, bug reports, ideas.

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Be kind, assume good intent, and treat others the way you'd like to be treated.

## Getting started

1. Fork the repository.
2. Create a branch: `git checkout -b feature/my-change` or `git checkout -b fix/my-bug`.
3. Make your change.
4. Add or update tests — every new behaviour should have one.
5. Run the checks (see below).
6. Open a pull request against `main`.

## Development environment

Pork is built with Go. You need Go 1.25+.

```bash
make build    # compile the binary into bin/pork
make test     # run all tests
make lint     # run golangci-lint
make vet      # run go vet
```

Prefer `make test` and `make lint` from the repository root — they pin the exact commands CI uses.

## What makes a good change

- **Keep it small.** A focused PR is easier to review and faster to merge.
- **Respect the existing style.** `gofmt` and `goimports` are enforced in CI.
- **Test the behaviour, not the implementation.** Tests should describe what Pork does, not how.
- **Document real usage.** If you change a command or its flags, update the README usage section.
- **Update the changelog.** Add a bullet under the next unreleased version in `CHANGELOG.md` following its format.

## Commit messages

Write clear, imperative commit messages that explain **why**:

```
fix: exclude system processes from kill safety checks
```

## Pull request checklist

- [ ] `make test` passes
- [ ] `make lint` passes
- [ ] README updated if user-facing behaviour changed
- [ ] CHANGELOG.md updated
- [ ] No unrelated changes

## Where to ask questions

Open a [GitHub discussion](https://github.com/abraira85/pork/discussions) or an issue if you are unsure about anything. It's better to ask than to guess.