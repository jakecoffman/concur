package concur

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

type Task interface {
	Process()
	Print()
}

type Factory interface {
	Make(line string) Task
}

func Run(factory Factory) {
	var wg sync.WaitGroup

	in := make(chan Task)

	wg.Add(1)
	go func() {
		defer func() {
			close(in)
			wg.Done()
		}()
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if line != "" {
				in <- factory.Make(line)
			}
		}
		if s.Err() != nil {
			log.Fatalf("Error reading STDIN: %s", s.Err())
		}
	}()

	out := make(chan Task)

	// TODO: Make goroutine limit configurable
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				t.Process()
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
		t.Print()
	}
}
