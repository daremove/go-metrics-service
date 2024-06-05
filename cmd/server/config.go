package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Endpoint        string `json:"address"`
	StoreInterval   int    `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	Restore         bool   `json:"restore"`
	Dsn             string `json:"database_dsn"`
	LogLevel        string
	SigningKey      string
	CryptoKey       string `json:"crypto_key"`
	TrustedSubnet   string `json:"trusted_subnet"`
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
		endpoint        string
		storeInterval   int
		fileStoragePath string
		restore         bool
		dsn             string
		logLevel        string
		signingKey      string
		cryptoKey       string
		configFile      string
		trustedSubnet   string
	)

	flag.StringVar(&endpoint, "a", "", "address and port to run server")
	flag.IntVar(&storeInterval, "i", 0, "interval of saving data to disk")
	flag.StringVar(&fileStoragePath, "f", "", "path of storage file")
	flag.BoolVar(&restore, "r", false, "should server restore data from storage file")
	flag.StringVar(&dsn, "d", "", "data source name for database connection")
	flag.StringVar(&signingKey, "k", "", "data signing key")
	flag.StringVar(&cryptoKey, "crypto-key", "", "path to the encryption key")
	flag.StringVar(&configFile, "c", "cmd/server/default_config.json", "path to the configuration file")
	flag.StringVar(&trustedSubnet, "t", "", "CIDR of agent")
	flag.Parse()

	if address := os.Getenv("ADDRESS"); address != "" {
		endpoint = address
	}

	if si := os.Getenv("STORE_INTERVAL"); si != "" {
		v, err := strconv.Atoi(si)

		if err != nil {
			log.Fatalf("STORE_INTERVAL couldn't parsed %s", err)
		}

		storeInterval = v
	}

	if fs := os.Getenv("FILE_STORAGE_PATH"); fs != "" {
		fileStoragePath = fs
	}

	if r := os.Getenv("RESTORE"); r != "" {
		v, err := strconv.ParseBool(r)

		if err != nil {
			log.Fatalf("RESTORE couldn't parsed %s", err)
		}

		restore = v
	}

	if d := os.Getenv("DATABASE_DSN"); d != "" {
		dsn = d
	}

	if l := os.Getenv("LOG_LEVEL"); l != "" {
		logLevel = "debug"
	} else {
		logLevel = "error"
	}

	if signingKeyEnv := os.Getenv("KEY"); signingKeyEnv != "" {
		signingKey = signingKeyEnv
	}

	if cryptoKeyEnv := os.Getenv("CRYPTO_KEY"); cryptoKeyEnv != "" {
		cryptoKey = cryptoKeyEnv
	}

	if configFileEnv := os.Getenv("CONFIG"); configFileEnv != "" {
		configFile = configFileEnv
	}

	if trustedSubnetEnv := os.Getenv("TRUSTED_SUBNET"); trustedSubnetEnv != "" {
		trustedSubnet = trustedSubnetEnv
	}

	if configFile != "" {
		fileConfig, err := loadConfigFromFile(configFile)

		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}

		if endpoint == "" {
			endpoint = fileConfig.Endpoint
		}

		if storeInterval == 0 {
			storeInterval = fileConfig.StoreInterval
		}

		if fileStoragePath == "" {
			fileStoragePath = fileConfig.FileStoragePath
		}

		if !restore {
			restore = fileConfig.Restore
		}

		if dsn == "" {
			dsn = fileConfig.Dsn
		}

		if cryptoKey == "" {
			cryptoKey = fileConfig.CryptoKey
		}

		if trustedSubnet == "" {
			trustedSubnet = fileConfig.TrustedSubnet
		}
	}

	return Config{
		endpoint,
		storeInterval,
		fileStoragePath,
		restore,
		dsn,
		logLevel,
		signingKey,
		cryptoKey,
		trustedSubnet,
	}
}
