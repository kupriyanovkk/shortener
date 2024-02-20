package config

import (
	"bytes"
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
	TrustedSubnet   string `json:"trusted_subnet"`
	ConfigFile      string
	EnableGRPC      bool
}

// ParseFlags parses and retrieves environment variables.
func ParseFlags(progname string, args []string) (*ConfigFlags, error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var (
		serverAddress   string
		baseURL         string
		fileStoragePath string
		databaseDSN     string
		enableHTTPS     bool
		configFile      string
		trustedSubnet   string
		enableGRPC      bool
	)

	parsedFlags := ConfigFlags{}

	flags.StringVar(&serverAddress, "a", "", "address and port to run server")
	flags.StringVar(&baseURL, "b", "", "the address of the resulting shortened URL")
	flags.StringVar(&fileStoragePath, "f", "", "the full name of the file where the data is saved in JSON")
	flags.StringVar(&databaseDSN, "d", "", "the address for DB connection")
	flags.BoolVar(&enableHTTPS, "s", false, "enable HTTPS support")
	flags.StringVar(&configFile, "c", "", "path to config file")
	flags.StringVar(&configFile, "config", "", "path to config file")
	flags.StringVar(&trustedSubnet, "t", "", "trusted subnet")
	flags.BoolVar(&enableGRPC, "g", false, "enable gRPC support")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	configFromEnv := os.Getenv("CONFIG")
	if configFromEnv != "" {
		configFile = configFromEnv
	}
	parsedFlags.ConfigFile = configFile
	parsedFlags.EnableHTTPS = enableHTTPS

	if configFile != "" {
		configData, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(configData, &parsedFlags)
		if err != nil {
			log.Fatal(err)
		}
	}

	updateIfNotEmpty := func(value, envValue string, field *string) {
		if envValue != "" {
			*field = envValue
		} else if value != "" {
			*field = value
		}
	}

	updateIfNotEmpty(serverAddress, os.Getenv("SERVER_ADDRESS"), &parsedFlags.ServerAddress)
	updateIfNotEmpty(baseURL, os.Getenv("BASE_URL"), &parsedFlags.BaseURL)
	updateIfNotEmpty(fileStoragePath, os.Getenv("FILE_STORAGE_PATH"), &parsedFlags.FileStoragePath)
	updateIfNotEmpty(databaseDSN, os.Getenv("DATABASE_DSN"), &parsedFlags.DatabaseDSN)
	updateIfNotEmpty(trustedSubnet, os.Getenv("TRUSTED_SUBNET"), &parsedFlags.TrustedSubnet)

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		parsedFlags.EnableHTTPS = envEnableHTTPS == "true"
	}

	if parsedFlags.ServerAddress == "" {
		parsedFlags.ServerAddress = "localhost:8080"
	}
	if parsedFlags.BaseURL == "" {
		parsedFlags.BaseURL = "http://localhost:8080"
	}

	return &parsedFlags, nil
}

// App structure contains flags, store and URLchan.
type App struct {
	Flags   *ConfigFlags
	Store   storeInterface.Store
	URLChan chan storeInterface.DeletedURLs
}
