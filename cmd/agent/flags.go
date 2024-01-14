package main

import (
	"flag"
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
}
