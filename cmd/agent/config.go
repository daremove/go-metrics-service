package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	endpoint       string
	reportInterval uint64
	pollInterval   uint64
	signingKey     string
}

func NewConfig() Config {
	var (
		endpoint       string
		reportInterval uint64
		pollInterval   uint64
		signingKey     string
	)

	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port where to send data")
	flag.Uint64Var(&reportInterval, "r", 10, "frequency of sending data to server")
	flag.Uint64Var(&pollInterval, "p", 2, "frequency of polling stats data")
	flag.StringVar(&signingKey, "k", "", "data signing key")
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

	return Config{
		endpoint,
		reportInterval,
		pollInterval,
		signingKey,
	}
}
