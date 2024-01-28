package config

import (
	"flag"
	"log"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		envVariables   map[string]string
		configFile     string
		expectedOutput ConfigFlags
	}{
		{
			name:         "JSONConfigFile",
			envVariables: map[string]string{},
			configFile:   "./test_config_file.json",
			expectedOutput: ConfigFlags{
				ServerAddress:   "json_server_address",
				BaseURL:         "json_base_url",
				FileStoragePath: "json_file_storage_path",
				DatabaseDSN:     "json_database_dsn",
				EnableHTTPS:     false,
			},
		},
		{
			name: "EnvironmentVariables",
			envVariables: map[string]string{
				"SERVER_ADDRESS":    "test_server_address",
				"BASE_URL":          "test_base_url",
				"FILE_STORAGE_PATH": "test_file_storage_path",
				"DATABASE_DSN":      "test_database_dsn",
				"ENABLE_HTTPS":      "true",
			},
			configFile: "./test_config_file.json",
			expectedOutput: ConfigFlags{
				ServerAddress:   "test_server_address",
				BaseURL:         "test_base_url",
				FileStoragePath: "test_file_storage_path",
				DatabaseDSN:     "test_database_dsn",
				EnableHTTPS:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVariables {
				os.Setenv(key, value)
			}

			// Create and write to config file
			if tt.configFile != "" {
				configData := `{"server_address":"json_server_address","base_url":"json_base_url","file_storage_path":"json_file_storage_path","database_dsn":"json_database_dsn","enable_https":false}`
				err := os.WriteFile(tt.configFile, []byte(configData), 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(tt.configFile)
				os.Setenv("CONFIG", tt.configFile)
			}

			flagSet := flag.NewFlagSet("test", flag.ExitOnError)
			parsedFlags := ParseFlags(flagSet)
			if parsedFlags != tt.expectedOutput {
				t.Errorf("Expected %v, but got %v", tt.expectedOutput, parsedFlags)
			}
		})
	}
}
