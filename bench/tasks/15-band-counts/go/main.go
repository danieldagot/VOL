package main

import "fmt"

func main() {
	vals := []int{22, 18, 25, 31, 29, 17, 24, 28, 20, 26}
	total, hot, mild, cold := 0, 0, 0, 0
	for _, v := range vals {
		total += v
		if v >= 28 {
			hot++
		} else if v >= 20 {
			mild++
		} else {
			cold++
		}
	}
	fmt.Println(len(vals))
	fmt.Println(total / len(vals))
	fmt.Println(hot)
	fmt.Println(mild)
	fmt.Println(cold)
}
