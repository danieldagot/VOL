package main

import "fmt"

func main() {
	nums := []int{3, 8, 1, 12, 5, 9, 4, 15, 2, 11}
	largeCount, largeSum, largeDouble, smallCount := 0, 0, 0, 0
	for _, n := range nums {
		if n > 5 {
			largeCount++
			largeSum += n
			largeDouble += n * 2
		}
		if n < 5 {
			smallCount++
		}
	}
	fmt.Println(largeCount)
	fmt.Println(largeSum)
	fmt.Println(largeDouble)
	fmt.Println(smallCount)
}
