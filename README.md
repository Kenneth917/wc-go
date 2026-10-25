# wc-go

A Go clone of `wc`, built test-first for the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc) series.

- `make unit`: TDD unit tests for the counting logic
- `make bdd`: Godog scenarios in `features/` against a locally built binary
- `make e2e`: the same scenarios against the Docker image

## Try it

```bash
go build -o wc ./cmd/wc
./wc examples/sample.txt        # lines, words, bytes
./wc -m examples/sample.txt     # characters
cat examples/sample.txt | ./wc -l
```

## Contributing

`main` is protected. Work on a feature branch and open a pull request; CI (`go vet`, unit specs, BDD scenarios, Docker e2e) must pass before it can merge.
