package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jakecoffman/concur"
)

func main() {
	concur.Run(&myFactory{})
}

type myFactory struct{}

func (t *myFactory) Make(line string) concur.Task {
	return &myTask{url: line}
}

type myTask struct {
	url       string
	err       error
	duration  float64
	bytesRead int64
}

func (t *myTask) Process() {
	start := time.Now()
	r, err := http.Get(t.url)
	if err != nil {
		t.err = err
		return
	}
	t.bytesRead, err = io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		t.err = err
		return
	}
	r.Body.Close()
	t.duration = time.Since(start).Seconds()
}

func (t myTask) Print() {
	if t.err != nil {
		fmt.Printf("%30v\t%v\n", t.url, t.err)
	} else {
		fmt.Printf("%30v\t%v\t%vs\n", t.url, t.bytesRead, t.duration)
	}
}
