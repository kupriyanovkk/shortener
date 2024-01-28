package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

// ConfigFlags contains flags for app.
type ConfigFlags struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	ConfigFile      string
}

// ParseFlags using for parsing and getting environment variables.
func ParseFlags(flag *flag.FlagSet) ConfigFlags {
	var (
		runAddress      string
		baseURL         string
		fileStoragePath string
		databaseDSN     string
		enableHTTPS     bool
		configFile      string
	)

	parsedFlags := ConfigFlags{}

	flag.StringVar(&runAddress, "a", "", "address and port to run server")
	flag.StringVar(&baseURL, "b", "", "the address of the resulting shortened URL")
	flag.StringVar(&fileStoragePath, "f", "", "the full name of the file where the data is saved in JSON")
	flag.StringVar(&databaseDSN, "d", "", "the address for DB connection")
	flag.BoolVar(&enableHTTPS, "s", false, "enable HTTPS support")
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.StringVar(&configFile, "config", "", "path to config file")

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		configFile = envConfig
	}

	if configFile != "" {
		file, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(file, &parsedFlags)
		if err != nil {
			log.Fatal(err)
		}
	}

	if runAddress != "" {
		parsedFlags.ServerAddress = runAddress
	}
	if baseURL != "" {
		parsedFlags.BaseURL = baseURL
	}
	if fileStoragePath != "" {
		parsedFlags.FileStoragePath = fileStoragePath
	}
	if databaseDSN != "" {
		parsedFlags.DatabaseDSN = databaseDSN
	}
	if enableHTTPS {
		parsedFlags.EnableHTTPS = enableHTTPS
	}

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		parsedFlags.ServerAddress = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		parsedFlags.BaseURL = envBaseAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		parsedFlags.FileStoragePath = envFileStoragePath
	}
	if envDatabaseDNS := os.Getenv("DATABASE_DSN"); envDatabaseDNS != "" {
		parsedFlags.DatabaseDSN = envDatabaseDNS
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		parsedFlags.EnableHTTPS = envEnableHTTPS == "true"
	}

	if parsedFlags.ServerAddress == "" {
		parsedFlags.ServerAddress = "localhost:8080"
	}
	if parsedFlags.BaseURL == "" {
		parsedFlags.BaseURL = "http://localhost:8080"
	}

	return parsedFlags
}

// App structure contains flags, store and URLchan.
type App struct {
	Flags   ConfigFlags
	Store   storeInterface.Store
	URLChan chan storeInterface.DeletedURLs
}
