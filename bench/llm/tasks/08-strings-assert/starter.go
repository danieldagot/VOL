package main

import "fmt"

func main() {
	word := "hello"
	length := len(word)
	doubled := word + "!"
	fmt.Println(length)
	fmt.Println(doubled)
	if length != 6 {
		panic("Expected length 5")
	}
}
