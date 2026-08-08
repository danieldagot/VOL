package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimSpace("  vol  "))
	fmt.Println(strings.Join(strings.Split("a,b,c", ","), "-"))
	fmt.Println(strings.Contains("vocabulary", "cab"))
	fmt.Println(strings.ReplaceAll("a-a-a", "-", "+"))
}
