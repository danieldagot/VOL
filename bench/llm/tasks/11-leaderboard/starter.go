package main

import "fmt"

func main() {
	player1 := []int{8, 6, 9, 5, 10, 7}
	player2 := []int{7, 8, 6, 9, 8, 8}
	p1Total, p2Total := 0, 0
	for _, score := range player1 {
		p1Total += score
	}
	for _, score := range player2 {
		p2Total += score
	}
	fmt.Println("Player 1 total: " + fmt.Sprint(p1Total))
	fmt.Println("Player 2 total: " + fmt.Sprint(p2Total))
}
