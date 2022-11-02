## Golang Proto Danksharding

## Intro

This is a golang version of the proto danksharding specs using gnark.

Audit of gnark code: <https://github.com/ConsenSys/gnark-crypto/blob/master/audit_oct2022.pdf>

In particular, we only use the group operations and pairings code.

## Test Vectors

To generate test vectors:

```
$ cd test_vectors
$ go run *.go
```

This will produce a series of json files.