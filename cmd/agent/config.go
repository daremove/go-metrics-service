package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Endpoint       string `json:"address"`
	ReportInterval uint64 `json:"report_interval"`
	PollInterval   uint64 `json:"poll_interval"`
	SigningKey     string
	RateLimit      uint64
	CryptoKey      string `json:"crypto_key"`
}

func loadConfigFromFile(path string) (Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func NewConfig() Config {
	var (
		endpoint       string
		reportInterval uint64
		pollInterval   uint64
		signingKey     string
		rateLimit      uint64
		cryptoKey      string
		configFile     string
	)

	flag.StringVar(&endpoint, "a", "", "address and port where to send data")
	flag.Uint64Var(&reportInterval, "r", 0, "frequency of sending data to server")
	flag.Uint64Var(&pollInterval, "p", 0, "frequency of polling stats data")
	flag.StringVar(&signingKey, "k", "", "data signing key")
	flag.Uint64Var(&rateLimit, "l", 1, "rate limit of batched request")
	flag.StringVar(&cryptoKey, "crypto-key", "", "path to the encryption key")
	flag.StringVar(&configFile, "c", "cmd/agent/default_config.json", "path to the configuration file")
	flag.Parse()

	if address := os.Getenv("ADDRESS"); address != "" {
		endpoint = address
	}

	if reportIntervalEnv := os.Getenv("REPORT_INTERVAL"); reportIntervalEnv != "" {
		value, err := strconv.Atoi(reportIntervalEnv)

		if err != nil {
			log.Fatal(err)
		}

		reportInterval = uint64(value)
	}

	if pollIntervalEnv := os.Getenv("POLL_INTERVAL"); pollIntervalEnv != "" {
		value, err := strconv.Atoi(pollIntervalEnv)

		if err != nil {
			log.Fatal(err)
		}

		pollInterval = uint64(value)
	}

	if signingKeyEnv := os.Getenv("KEY"); signingKeyEnv != "" {
		signingKey = signingKeyEnv
	}

	if rateLimitEnv := os.Getenv("RATE_LIMIT"); rateLimitEnv != "" {
		value, err := strconv.Atoi(rateLimitEnv)

		if err != nil {
			log.Fatal(err)
		}

		rateLimit = uint64(value)
	}

	if cryptoKeyEnv := os.Getenv("CRYPTO_KEY"); cryptoKeyEnv != "" {
		cryptoKey = cryptoKeyEnv
	}

	if configFileEnv := os.Getenv("CONFIG"); configFileEnv != "" {
		configFile = configFileEnv
	}

	if configFile != "" {
		fileConfig, err := loadConfigFromFile(configFile)

		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}

		if endpoint == "" {
			endpoint = fileConfig.Endpoint
		}

		if reportInterval == 0 {
			reportInterval = fileConfig.ReportInterval
		}

		if pollInterval == 0 {
			pollInterval = fileConfig.PollInterval
		}

		if cryptoKey == "" {
			cryptoKey = fileConfig.CryptoKey
		}
	}

	return Config{
		endpoint,
		reportInterval,
		pollInterval,
		signingKey,
		rateLimit,
		cryptoKey,
	}
}
