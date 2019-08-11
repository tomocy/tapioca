package client

import (
	"os"

	"github.com/tomocy/tapioca/domain"
)

func newCLI(cnf config) *CLI {
	return &CLI{
		cnf:     cnf,
		printer: newPrinter(cnf.format),
	}
}

type CLI struct {
	cnf     config
	printer printer
}

func (c *CLI) Run() error {
	if c.cnf.author != "" {
		return c.summarizeAuthorCommitsOfToday()
	}

	return c.summarizeCommitsOfToday()
}

func (c *CLI) summarizeAuthorCommitsOfToday() error {
	s, err := summarizeAuthorCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name, c.cnf.author)
	if err != nil {
		return err
	}

	c.showSummary(*s)

	return nil
}

func (c *CLI) summarizeCommitsOfToday() error {
	s, err := summarizeCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name)
	if err != nil {
		return err
	}

	c.showSummary(*s)

	return nil
}

func (c *CLI) showSummary(s domain.Summary) {
	c.printer.PrintSummary(os.Stdout, s)
}
