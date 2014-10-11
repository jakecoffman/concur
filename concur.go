package concur

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Task interface {
	process()
	print()
}

type Factory interface {
	make(line string) Task
}

type myFactory struct{}

func (t *myFactory) make(line string) Task {
	return &myTask{url: line}
}

type myTask struct {
	url       string
	err       error
	duration  float64
	bytesRead int64
}

func (t *myTask) process() {
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

func (t myTask) print() {
	if t.err != nil {
		fmt.Printf("%v: (ERROR) %v\n", t.url, t.err)
	} else {
		fmt.Printf("%v: bytes: %v duration: %vs\n", t.url, t.bytesRead, t.duration)
	}
}

func run(f Factory) {
	var wg sync.WaitGroup

	in := make(chan Task)

	wg.Add(1)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			in <- f.make(s.Text())
		}
		if s.Err() != nil {
			log.Fatalf("Error reading STDIN: %s", s.Err())
		}
		close(in)
		wg.Done()
	}()

	out := make(chan Task)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				t.process()
				out <- t
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for t := range out {
		t.print()
	}
}

func main() {
	run(&myFactory{})
}
