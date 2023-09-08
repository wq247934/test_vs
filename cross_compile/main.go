package main

import (
	"flag"
	"test_vs/cross_compile/hello"
)

func main() {
	flag.Parse()

	name := flag.Arg(0)
	if name == "" {
		name = "stranger"
	}

	hello.Greet(name)
}
