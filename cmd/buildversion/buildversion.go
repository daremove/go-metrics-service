// Package buildversion предназначен для отображения текущей версии сборки.
package buildversion

import "fmt"

var (
	BuildVersion = "NA"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func printVersion() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}

func init() {
	printVersion()
}
