package main

import (
	"fmt"
	"time"
)

/*
Написать 3 функции:
writer - генерит числа от 1 до 10
doubler - умножает числа на 2, имитируя работу (500ms)
reader - читает и выводит на экран
*/

func main() {
	reader(double(writer()))
}

func reader(in <-chan int) {
	for data := range in {
		fmt.Println(data)
	}
}

func double(ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for data := range ch {
			out <- data * 2
			time.Sleep(500 * time.Millisecond)
		}
		close(out)
	}()
	return out
}

func writer() <-chan int {
	ch := make(chan int)

	go func() {
		for i := range 10 {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
