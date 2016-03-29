#### TCP Drain Test

This is an attempt to confirm or deny [this issue](https://github.com/google/go-github/pull/317) with closing http response bodies.

These are my test results so far:

```
$ go version
go version go1.6 darwin/amd64

$ PORT=5099 go run main.go 10
Without drain (x10): 2.989648ms
   With drain (x10): 784.782µs

$ PORT=5099 go run main.go 100
Without drain (x100): 12.047307ms
   With drain (x100): 6.030919ms

$ PORT=5099 go run main.go 1000
Without drain (x1000): 65.449043ms
   With drain (x1000): 55.154988ms

$ PORT=5099 go run main.go 10000
Without drain (x10000): 581.107868ms
   With drain (x10000): 571.955831ms

$ PORT=5099 go run main.go 100000
Without drain (x100000): 5.692069876s
   With drain (x100000): 5.648899366s
```