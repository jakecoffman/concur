concur
======

This package is an easy way to make concurrent Unix-style line-by-line executables.

First, implement a Factory that will receive a single line of input and return a
Task that will be executed concurrently. 

```go
type Factory interface {
	Make(line string) Task
}
```

You also need to implement a Task:

```go
type Task interface {
	Process()
	Print()
}
```

Process() is where your actual work will be done. Print() will be called once
Process() is complete.

example
-------

The example directory contains code for an executable that takes a list of urls (input.txt)
and retrieves them, counting the number of bytes and the time it took to perform the request.

To run it:
```sh
go get github.com/jakecoffman/concur/example
cd $GOPATH/src/github.com/jakecoffman/concur/example
cat input.txt | go run main.go
```

You should see output similar to this:

```
             http://golang.org     7261 0.136367914s
               http://perl.org    13873 0.296372999s
            http://clojure.org    35902 0.769628144s
             http://python.org    45659 0.771808377s
          http://ruby-lang.org      833 0.908750531s
            http://haskell.org    21424 1.005203749s
          http://rust-lang.org     9544 1.56502838s
```

Running it again may produce results in a different order since this is concurrent!
