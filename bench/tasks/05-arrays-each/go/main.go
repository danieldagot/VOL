package main

import "fmt"

func main() {
	scores := []int{72, 95, 81, 64}
	fmt.Println("Students: " + fmt.Sprint(len(scores)))
	scores[3] = 70
	for _, score := range scores {
		if score >= 80 {
			fmt.Println("High score: " + fmt.Sprint(score))
		}
	}
}
