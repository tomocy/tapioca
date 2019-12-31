package main

import (
	"flag"
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

type help struct {
	err error
}

func (h help) Run() error {
	flag.Usage()
	return h.err
}

type config struct {
	owner        string
	repo         string
	since, until time.Time
}

type presenter interface {
	PresentSummary(domain.Summary)
}
