package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeJSONFile(t *testing.T, path string, content interface{}) {
	file, err := os.Create(path)
	require.NoError(t, err)

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)

	require.NoError(t, err)
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Run("Should load valid config from JSON file", func(t *testing.T) {
		configContent := Config{
			Endpoint:       "localhost:8080",
			ReportInterval: 10,
			PollInterval:   2,
			CryptoKey:      "path/to/crypto_key.pem",
		}
		configPath := "test_config.json"

		writeJSONFile(t, configPath, configContent)
		defer os.Remove(configPath)

		config, err := loadConfigFromFile(configPath)

		require.NoError(t, err)
		assert.Equal(t, configContent, config)
	})

	t.Run("Should return error for invalid JSON file", func(t *testing.T) {
		invalidContent := []byte("invalid json content")
		invalidPath := "invalid_config.json"

		err := os.WriteFile(invalidPath, invalidContent, 0644)
		defer os.Remove(invalidPath)

		require.NoError(t, err)

		_, err = loadConfigFromFile(invalidPath)

		assert.Error(t, err)
	})

	t.Run("Should return error for non-existent file", func(t *testing.T) {
		nonExistentPath := "non_existent_config.json"

		_, err := loadConfigFromFile(nonExistentPath)

		assert.Error(t, err)
	})
}
