package main

import "fmt"

func main() {
	temperature := 24
	sunny := true
	if temperature >= 20 && sunny {
		fmt.Println("Good weather for a walk.")
	} else {
		fmt.Println("Stay inside and write VOL.")
	}
}
