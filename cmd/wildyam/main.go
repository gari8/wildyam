package main

import (
	"flag"
	"fmt"
)

type Recepter struct {
	SubCmd int
}

func main() {
	var recepter Recepter
	flag.Parse()
	fmt.Println(flag.Args())
	fmt.Println(recepter)
}
