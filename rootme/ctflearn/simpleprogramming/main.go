package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

func main() {
	file, _ := os.Open("data.dat")
	var ones, zeros, count int64
	for bin := range generator(file) {
		binrune := []rune(bin)
		atomic.CompareAndSwapInt64(&ones, ones, 0)
		atomic.CompareAndSwapInt64(&zeros, zeros, 0)
		for i := 0; i < len(binrune); i++ {
			switch binrune[i] {
			case '0':
				atomic.AddInt64(&zeros, 1)
			case '1':
				atomic.AddInt64(&ones, 1)
			}
		}
		if (zeros%3 == 0) || (ones%2 == 0) {
			atomic.AddInt64(&count, 1)

		}
	}
	fmt.Printf("Total count = %d\n", count)
}

func generator(r io.Reader) <-chan string {
	binary := make(chan string)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	go func() {
		defer close(binary)
		for scanner.Scan() {
			binary <- scanner.Text()
		}
	}()
	return binary
}
