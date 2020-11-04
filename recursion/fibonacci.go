package main

import "fmt"

func main() {
	fmt.Println(fib(19))
}

// find the nth number in the fibonacci sequence
func fib(n int) int {
	if n <= 2 {
		return 1
	}

	return fib(n - 1) + fib(n - 2)
}