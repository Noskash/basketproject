package main

import (
	"fmt"

	"github.com/Noskash/basketproject/internal/src"
)

func main() {
	res := src.Get_json_file("56629288")
	fmt.Printf(res)
}
