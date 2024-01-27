package config

import (
	"flag"
	"os"

	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

// ConfigFlags contains flags for app.
type ConfigFlags struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	EnableHTTPS     bool
}

// ParseFlags using for parsing and getting environment variables.
func ParseFlags() ConfigFlags {
	var (
		runAddress      string
		baseURL         string
		fileStoragePath string
		databaseDSN     string
		enableHTTPS     bool
	)

	flag.StringVar(&runAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "the address of the resulting shortened URL")
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "the full name of the file where the data is saved in JSON")
	flag.StringVar(&databaseDSN, "d", "", "the address for DB connection")
	flag.BoolVar(&enableHTTPS, "s", false, "enable HTTPS support")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		runAddress = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		baseURL = envBaseAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envDatabaseDNS := os.Getenv("DATABASE_DSN"); envDatabaseDNS != "" {
		databaseDSN = envDatabaseDNS
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		enableHTTPS = envEnableHTTPS == "true"
	}

	return ConfigFlags{
		ServerAddress:   runAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		DatabaseDSN:     databaseDSN,
		EnableHTTPS:     enableHTTPS,
	}
}

// App structure contains flags, store and URLchan.
type App struct {
	Flags   ConfigFlags
	Store   storeInterface.Store
	URLChan chan storeInterface.DeletedURLs
}
