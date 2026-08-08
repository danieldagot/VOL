package main

import "fmt"

func main() {
	xs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	midCount, highSum := 0, 0
	for _, x := range xs {
		n := x * 3
		if n > 10 && n < 20 {
			midCount++
		}
		if n > 10 {
			highSum += n
		}
	}
	fmt.Println(midCount)
	fmt.Println(highSum)
}
