package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_Args(t *testing.T) {
	os.Clearenv()

	flags, _ := ParseFlags(os.Args[0], []string{"-a", "test_server_address", "-b", "test_base_url", "-f", "test_file_storage_path", "-d", "test_database_dsn", "-s", "true"})

	assert.Equal(t, "test_server_address", flags.ServerAddress, "ServerAddress not parsed correctly")
	assert.Equal(t, "test_base_url", flags.BaseURL, "BaseURL not parsed correctly")
	assert.Equal(t, "test_file_storage_path", flags.FileStoragePath, "FileStoragePath not parsed correctly")
	assert.Equal(t, "test_database_dsn", flags.DatabaseDSN, "DatabaseDSN not parsed correctly")
	assert.Equal(t, true, flags.EnableHTTPS, "EnableHTTPS not parsed correctly")
}

func TestParseFlags_JSON(t *testing.T) {
	os.Clearenv()

	configData := `{"server_address":"json_server_address","base_url":"json_base_url","file_storage_path":"json_file_storage_path","database_dsn":"json_database_dsn","enable_https":false}`
	configFile := "./test_config_file.json"
	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(configFile)

	os.Setenv("CONFIG", "./test_config_file.json")

	flags, _ := ParseFlags(os.Args[0], []string{})

	assert.Equal(t, "json_server_address", flags.ServerAddress, "ServerAddress not parsed correctly")
	assert.Equal(t, "json_base_url", flags.BaseURL, "BaseURL not parsed correctly")
	assert.Equal(t, "json_file_storage_path", flags.FileStoragePath, "FileStoragePath not parsed correctly")
	assert.Equal(t, "json_database_dsn", flags.DatabaseDSN, "DatabaseDSN not parsed correctly")
	assert.Equal(t, false, flags.EnableHTTPS, "EnableHTTPS not parsed correctly")
	assert.Equal(t, configFile, flags.ConfigFile, "ConfigFile not parsed correctly")
}

func TestParseFlags_ENV(t *testing.T) {
	os.Clearenv()

	configData := `{"server_address":"json_server_address","base_url":"json_base_url","file_storage_path":"json_file_storage_path","database_dsn":"json_database_dsn","enable_https":false}`
	configFile := "./test_config_file.json"
	err := os.WriteFile(configFile, []byte(configData), 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(configFile)

	os.Setenv("CONFIG", "./test_config_file.json")

	os.Setenv("SERVER_ADDRESS", "env_server_address")
	os.Setenv("BASE_URL", "env_base_url")
	os.Setenv("FILE_STORAGE_PATH", "env_file_storage_path")
	os.Setenv("DATABASE_DSN", "env_database_dsn")
	os.Setenv("ENABLE_HTTPS", "true")

	flags, _ := ParseFlags(os.Args[0], []string{"-a", "test_server_address", "-b", "test_base_url", "-f", "test_file_storage_path", "-d", "test_database_dsn", "-s", "true"})

	assert.Equal(t, "env_server_address", flags.ServerAddress, "ServerAddress not parsed correctly")
	assert.Equal(t, "env_base_url", flags.BaseURL, "BaseURL not parsed correctly")
	assert.Equal(t, "env_file_storage_path", flags.FileStoragePath, "FileStoragePath not parsed correctly")
	assert.Equal(t, "env_database_dsn", flags.DatabaseDSN, "DatabaseDSN not parsed correctly")
	assert.Equal(t, true, flags.EnableHTTPS, "EnableHTTPS not parsed correctly")
}
