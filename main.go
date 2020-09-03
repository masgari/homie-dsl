package main

import (
	"os"
	"github.com/masgari/homie-dsl/dsl"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: homie-dsl <config file path>")
	}
	dsl.Run(os.Args[1])
}
