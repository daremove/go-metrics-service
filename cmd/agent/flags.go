package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	endpoint       string
	reportInterval uint64
	pollInterval   uint64
)

func parseFlags() {
	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port where to send data")
	flag.Uint64Var(&reportInterval, "r", 10, "frequency of sending data to server")
	flag.Uint64Var(&pollInterval, "p", 2, "frequency of polling stats data")
	flag.Parse()

	if address := os.Getenv("ADDRESS"); address != "" {
		endpoint = address
	}

	if reportIntervalEnv := os.Getenv("REPORT_INTERVAL"); reportIntervalEnv != "" {
		value, err := strconv.Atoi(reportIntervalEnv)

		if err != nil {
			panic(err)
		}

		reportInterval = uint64(value)
	}

	if pollIntervalEnv := os.Getenv("POLL_INTERVAL"); pollIntervalEnv != "" {
		value, err := strconv.Atoi(pollIntervalEnv)

		if err != nil {
			panic(err)
		}

		pollInterval = uint64(value)
	}
}
