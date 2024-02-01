package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	endpoint        string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func NewConfig() Config {
	var (
		endpoint        string
		storeInterval   int
		fileStoragePath string
		restore         bool
	)

	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&storeInterval, "i", 300, "interval of saving data to disk")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "path of storage file")
	flag.BoolVar(&restore, "r", true, "should server restore data from storage file")
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

	return Config{endpoint, storeInterval, fileStoragePath, restore}
}
