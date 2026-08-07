package main

import "fmt"

func main() {
	temps := []int{22, 18, 25, 31, 29, 17, 24, 28, 20, 26}
	total, hot, mild, cold := 0, 0, 0, 0
	for _, t := range temps {
		total += t
		if t >= 28 {
			hot++
		} else if t >= 20 {
			mild++
		} else {
			cold++
		}
	}
	avg := total / len(temps)
	if hot+mild+cold != len(temps) || avg != 24 {
		panic("invalid temperature report")
	}
	fmt.Println("Days measured: " + fmt.Sprint(len(temps)))
	fmt.Println("Average: " + fmt.Sprint(avg))
	fmt.Println("Hot days (28+): " + fmt.Sprint(hot))
	fmt.Println("Mild days: " + fmt.Sprint(mild))
	fmt.Println("Cold days (<20): " + fmt.Sprint(cold))
}
