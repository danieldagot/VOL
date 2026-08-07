package main

import "fmt"

func main() {
	word := "hello"
	length := len(word)
	doubled := word + word
	fmt.Println(length)
	fmt.Println(doubled)
	if length != 5 {
		panic("Expected length 5")
	}
}
