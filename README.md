# wc-go

A Go clone of `wc`, built test-first for the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc) series.

- `make unit`: TDD unit tests for the counting logic
- `make bdd`: Godog scenarios in `features/` against a locally built binary
- `make e2e`: the same scenarios against the Docker image
