package config

import (
	"flag"
	"os"

	"github.com/kupriyanovkk/shortener/internal/storage"
)

type ConfigFlags struct {
	A string
	B string
	F string
	D string
}

func ParseFlags() ConfigFlags {
	var a string
	var b string
	var f string
	var d string

	flag.StringVar(&a, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&b, "b", "http://localhost:8080", "the address of the resulting shortened URL")
	flag.StringVar(&f, "f", "/tmp/short-url-db.json", "the full name of the file where the data is saved in JSON")
	flag.StringVar(&d, "d", "", "the address for DB connection")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		a = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		b = envBaseAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		f = envFileStoragePath
	}
	if envDatabaseDNS := os.Getenv("DATABASE_DSN"); envDatabaseDNS != "" {
		d = envDatabaseDNS
	}

	return ConfigFlags{
		A: a,
		B: b,
		F: f,
		D: d,
	}
}

type Env struct {
	Flags   ConfigFlags
	Storage storage.StorageModel
}
