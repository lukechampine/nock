nock
-----

[![GoDoc](https://godoc.org/github.com/lukechampine/nock?status.svg)](https://godoc.org/github.com/lukechampine/nock)
[![Go Report Card](http://goreportcard.com/badge/github.com/lukechampine/nock)](https://goreportcard.com/report/github.com/lukechampine/nock)

```
go get github.com/lukechampine/nock
```

`nock` implements a simple Nock interpreter, according the spec available at
https://urbit.org/docs/nock/definition.

This interpreter assumes that its input is well-formed, and does not
support atoms larger than a machine `int`.

### Example ###

```go
decrement := nock.Parse(`[42 [8 [1 0] 8 [1 6 [5 [0 7] 4 0 6]
            [0 6] 9 2 [0 2] [4 0 6] 0 7] 9 2 0 1]]`)
program := nock.Cell(nock.Atom(42), decrement)
result := nock.Nock(program)
println(result) // 41
```