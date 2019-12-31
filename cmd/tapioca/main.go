package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tomocy/tapioca/cmd/tapioca/client"
	"github.com/tomocy/tapioca/domain"
)

func main() {
	c := client.New()
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}

type config struct {
	owner        string
	repo         string
	since, until time.Time
}

type presenter interface {
	PresentSummary(domain.Summary)
}
