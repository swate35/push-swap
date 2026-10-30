# push-swap

Two Go programs for the push-swap stack-sorting exercise:

- `checker` validates and executes instructions read from standard input.
- `push-swap` prints instructions that sort the supplied integers.

Build both executables from the repository root:

```sh
go build -o checker ./cmd/checker
go build -o push-swap ./cmd/push-swap
```

On Windows, use output names such as `checker.exe` and `push-swap.exe`.
Arguments can be supplied as one space-separated string or as separate
integers:

```sh
push-swap "4 67 3 87 23" | checker "4 67 3 87 23"
```

The solver uses breadth-first search to produce a shortest instruction
sequence for stacks of up to eight values. For larger stacks it uses binary
radix sorting over coordinate-compressed values, keeping the operation count
scalable. Invalid integers, duplicates, and invalid checker instructions print
`Error` to standard error.