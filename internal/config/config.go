package config

import (
	"flag"
	"os"
)

type ConfigFlags struct {
	A string
	B string
}

func ParseFlags() ConfigFlags {
	var a string
	var b string

	flag.StringVar(&a, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&b, "b", "http://localhost:8080", "the address of the resulting shortened URL")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		a = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		b = envBaseAddr
	}

	return ConfigFlags{
		A: a,
		B: b,
	}
}
