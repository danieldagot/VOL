package main

import "fmt"

func main() {
	revenue := []int{240, 175, 289, 150, 225, 199, 180, 178}
	total, highSum, highCount, budgetCount := 0, 0, 0, 0
	for _, r := range revenue {
		total += r
		if r >= 200 {
			highCount++
			highSum += r
		} else {
			budgetCount++
		}
	}
	fmt.Println("Total revenue: " + fmt.Sprint(total))
	fmt.Println("Premium orders (200+): " + fmt.Sprint(highCount))
	fmt.Println("Premium revenue: " + fmt.Sprint(highSum))
	fmt.Println("Budget orders: " + fmt.Sprint(budgetCount))
}
