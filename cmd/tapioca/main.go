package main

import (
	"fmt"
	"os"

	"github.com/tomocy/tapioca/cmd/tapioca/client"
)

func main() {
	c := client.New()
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}
