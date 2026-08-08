package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.ToSlash(filepath.Join("reports", "2026", "q1.json")))
	fmt.Println(filepath.Base("reports/2026/q1.json"))
	fmt.Println(filepath.ToSlash(filepath.Dir("reports/2026/q1.json")))
	fmt.Println(filepath.Ext("reports/2026/q1.json"))
}
