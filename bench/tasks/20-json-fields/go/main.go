package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var v map[string]any
	if err := json.Unmarshal([]byte(`{"n":3,"name":"vol"}`), &v); err != nil {
		panic(err)
	}
	fmt.Println(int(v["n"].(float64)))
	fmt.Println(v["name"].(string))
	out, err := json.Marshal(map[string]int{"n": 3})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
