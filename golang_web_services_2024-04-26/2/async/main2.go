package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func processData(val int) int {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return val * 2
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 100 {
			in <- i
		}
		close(in)
	}()

	now := time.Now()
	go processParallel(in, out, 5)

	for val := range out {
		fmt.Println(val)
	}
	fmt.Println(time.Since(now))
}

// операция должна выполняться не более 5 секунд
func processParallel(in <-chan int, out chan<- int, numWorkers int) {
	wg := &sync.WaitGroup{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range in {
				resultChan := make(chan int, 1)
				go func() {
					resultChan <- processData(data)
				}()
				select {
				case result := <-resultChan:
					out <- result
				case <-ctx.Done():
					return
				}
			}

		}()
	}
	wg.Wait()
	close(out)
}
