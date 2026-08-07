package main

import "fmt"

func main() {
	player1 := []int{8, 6, 9, 5, 10, 7}
	player2 := []int{7, 8, 6, 9, 8, 8}
	p1Total, p2Total := 0, 0
	for _, s := range player1 {
		p1Total += s
	}
	for _, s := range player2 {
		p2Total += s
	}
	winner := "Player 2"
	if p1Total > p2Total {
		winner = "Player 1"
	}
	p1Strong, p2Strong := 0, 0
	for _, s := range player1 {
		if s >= 8 {
			p1Strong++
		}
	}
	for _, s := range player2 {
		if s >= 8 {
			p2Strong++
		}
	}
	fmt.Println("Player 1 total: " + fmt.Sprint(p1Total))
	fmt.Println("Player 2 total: " + fmt.Sprint(p2Total))
	fmt.Println("Winner: " + winner)
	fmt.Println("P1 rounds 8+: " + fmt.Sprint(p1Strong))
	fmt.Println("P2 rounds 8+: " + fmt.Sprint(p2Strong))
}
