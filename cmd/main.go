package main

import (
	"log"
	"rapid-bridge/cmd/cli"
)

func main() {

	if err := cli.RootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
