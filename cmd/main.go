package main

import (
	"encoding/json"
	"fmt"

	"github.com/Noskash/basketproject/internal/src"
)

func main() {
	res := src.Get_json_file("56615507")
	data, _ := json.MarshalIndent(res, " ", " ")
	fmt.Printf(string(data))
}
