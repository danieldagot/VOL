package main

import (
	"fmt"
	"os"
)

func main() {
	_ = os.Setenv("VOL_DENSITY_PORT", "9090")
	if v, ok := os.LookupEnv("VOL_DENSITY_PORT"); ok {
		fmt.Println(v)
	} else {
		fmt.Println("missing")
	}
	if v, ok := os.LookupEnv("VOL_DENSITY_NO_SUCH_KEY"); ok {
		fmt.Println(v)
	} else {
		fmt.Println("fallback")
	}
}
