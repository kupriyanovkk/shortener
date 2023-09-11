package config

import (
	"flag"
)

type ConfigFlags struct {
	A string
	B string
}

func ParseFlags() ConfigFlags {
	var a string
	var b string

	flag.StringVar(&a, "a", "http://localhost:8080", "address and port to run server")
	flag.StringVar(&b, "b", "http://localhost:8080", "the address of the resulting shortened URL")
	flag.Parse()

	return ConfigFlags{
		A: a,
		B: b,
	}
}
