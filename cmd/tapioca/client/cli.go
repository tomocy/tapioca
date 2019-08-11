package client

import (
	"fmt"
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
	return c.summarizeCommitsOfToday()
}

func (c *CLI) summarizeCommitsOfToday() error {
	s, err := summarizeCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name)
	if err != nil {
		return err
	}

	c.showSummary(*s)

	return nil
}

func (c *CLI) fetchCommits() error {
	report := reportFunc("fetch commits")

	uc := newCommitUsecase()
	cs, err := uc.FetchCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name)
	if err != nil {
		return report(err)
	}

	c.showCommits(cs)

	return nil
}

func (c *CLI) showSummary(s domain.Summary) {
	c.printer.PrintSummary(os.Stdout, s)
}

func (*CLI) showCommits(cs domain.Commits) {
	for _, c := range cs {
		fmt.Println(c)
	}
}
