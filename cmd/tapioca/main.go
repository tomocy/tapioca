package main

import (
	"fmt"
	"os"

	"github.com/tomocy/tapioca/cmd/tapioca/client"
)

func main() {
	runner := new(client.CLI)
	if err := runner.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}
