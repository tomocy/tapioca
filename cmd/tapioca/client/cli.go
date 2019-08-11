package client

import (
	"os"

	"github.com/tomocy/tapioca/domain"
)

func newCLI(printer printer) *CLI {
	return &CLI{
		printer: printer,
	}
}

type CLI struct {
	printer printer
}

func (c *CLI) ShowSummary(s domain.Summary) {
	c.printer.PrintSummary(os.Stdout, s)
}
