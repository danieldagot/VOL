package main

import "fmt"

func main() {
	fmt.Println("Countdown")
	remaining := 3
	for remaining > 0 {
		fmt.Println(remaining)
		remaining--
	}
	for i := 0; i < 2; i++ {
		fmt.Println("Go!")
	}
}
