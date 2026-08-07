package main

import "fmt"

func main() {
	scores := []int{85, 72, 91, 60, 78, 95, 55, 68, 88, 74}
	total, aGrades, bGrades, passing, failing := 0, 0, 0, 0, 0
	for _, s := range scores {
		total += s
		if s >= 90 {
			aGrades++
		}
		if s >= 80 && s < 90 {
			bGrades++
		}
		if s >= 60 {
			passing++
		}
		if s < 60 {
			failing++
		}
	}
	avg := total / len(scores)
	fmt.Println("Class average: " + fmt.Sprint(avg))
	fmt.Println("A grades: " + fmt.Sprint(aGrades))
	fmt.Println("B grades: " + fmt.Sprint(bGrades))
	fmt.Println("Passing: " + fmt.Sprint(passing))
	fmt.Println("Failing: " + fmt.Sprint(failing))
}
