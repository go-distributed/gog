gog
===

A gossip-based broadcast system

[![GoDoc] (https://godoc.org/github.com/go-distributed/gog?status.png)](https://godoc.org/github.com/go-distributed/gog)

[Proposal](https://docs.google.com/document/d/1ouAsRyMZHtBKkpv4XbAD4vDEiH1wKa2lJ7AYm6F56ME/edit?usp=sharing)

[Design Document](https://docs.google.com/document/d/189erD25i-CLiYEYWVIo9OL6MKID2RnIKbFsByrTznkA/edit?usp=sharing)

###How to build:

```shell
$ go build
```

Run a standalone node:

```shell
$ ./gog -addr="localhost:8424"
```

To join an existing cluster:

```shell
$ ./gog -addr="localhost:8425" -join_node="localhost:8424"
```
