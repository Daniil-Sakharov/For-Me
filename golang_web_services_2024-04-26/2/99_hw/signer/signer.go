package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var in chan interface{}
	wg := &sync.WaitGroup{}
	for _, j := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(j, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			if d, ok := data.(string); ok {
				defer wg.Done()
				mu.Lock()
				md5 := DataSignerMd5(d)
				mu.Unlock()
				crc32 := DataSignerCrc32(d)
				secondCrc32 := DataSignerCrc32(md5)
				result := crc32 + "~" + secondCrc32
				out <- result
			}
		}(data)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	for data := range in {
		wg := &sync.WaitGroup{}
		slice := make([]string, 6)
		if d, ok := data.(string); ok {
			for i := 0; i < 6; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					th := strconv.Itoa(i)
					slice[i] = DataSignerCrc32(th + d)
				}(i)
			}
			wg.Wait()
			var result string
			for _, d := range slice {
				result += d
			}
			out <- result
		}
	}
}

func CombineResults(in, out chan interface{}) {
	allResult := []string{}
	for data := range in {
		if d, ok := data.(string); ok {
			allResult = append(allResult, d)
		}
	}
	sort.Strings(allResult)
	result := strings.Join(allResult, "_")
	out <- result
}
