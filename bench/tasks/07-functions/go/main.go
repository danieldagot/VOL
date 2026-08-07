package main

import "fmt"

func greet(name string) string {
	return "Hello, " + name
}

func square(number int) int {
	return number * number
}

func main() {
	fmt.Println(greet("friend"))
	fmt.Println("Six squared is " + fmt.Sprint(square(6)))
}
