package main

import "fmt"

func main() {
	numbers := []int{4, 7, 2, 9, 12}
	var large []int
	for _, n := range numbers {
		if n > 5 {
			large = append(large, n)
		}
	}
	total := 0
	for _, n := range large {
		total += n
	}
	fmt.Println("Sum: " + fmt.Sprint(total))
	if total != 28 {
		panic("Unexpected collection total.")
	}
}
