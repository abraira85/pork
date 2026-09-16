# Security Policy

## Supported versions

| Version  | Supported          |
| -------- | ------------------ |
| latest   | :white_check_mark: |
| < latest | :x:                |

Pork is under active development. Security fixes land on the latest release and are
tagged as patch releases (e.g. `v0.1.1`).

## Reporting a vulnerability

Pork reads process information from your local system and can kill processes you
ask it to. A vulnerability could mean an attacker tricking you into running Pork
against an endpoint they control, or the tool misbehaving when it encounters
malformed input.

If you believe you have found a security issue:

1. **Do not open a public issue.**
2. Email the maintainer directly at **<rober@outboss.io>** with as much context as
   you can safely share:
   - How the issue manifests (paste sanitized output if relevant),
   - What versions are affected,
   - What you expected to happen vs. what happened.

You will receive a response as soon as possible. Please allow time for a fix and a
coordinated disclosure before posting details publicly.

## Disclosure policy

- Issues are triaged privately and we work on a fix before announcing.
- We publish CVE or GHSA identifiers when appropriate.
- After a fix is released and announced, the issue may be discussed publicly.

## Safe harbour

Researchers who act in good faith are welcome. We will not pursue legal action for
responsible, non-destructive testing of this project.
