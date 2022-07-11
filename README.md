# A local command line

The basic idea is to use the idea of `Stack` and `LinkedList`.

When the command line is started, an initial `transaction` will be created, and all 
operations are based on this transaction. When a new transaction is created using `BEGIN`, will push the
new transaction on the top of the stack. If the `ROLLBACK` is performed, the current stack top will be discarded directly, 
while the `COMMIT` operation will overwrite all the values at the top of the stack to the previous transaction.

## How to build

You can use `go run cmd/*` directly.

```go
go build -o tkvs cmd/*
./tkvs
```
