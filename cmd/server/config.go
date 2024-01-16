package main

import (
	"flag"
	"os"
)

type Config struct {
	endpoint string
}

func NewConfig() Config {
	var endpoint string

	flag.StringVar(&endpoint, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if address := os.Getenv("ADDRESS"); address != "" {
		endpoint = address
	}

	return Config{endpoint}
}
