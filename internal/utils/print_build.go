package utils

import "fmt"

var (
	BuildVersion = "NA"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func init() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}
