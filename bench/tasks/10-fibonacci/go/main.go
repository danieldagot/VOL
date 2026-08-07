package main

import "fmt"

func main() {
	a, b := 0, 1
	for i := 0; i < 8; i++ {
		fmt.Println(a)
		a, b = b, a+b
	}
}
