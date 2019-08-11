package client

import (
	"fmt"

	"github.com/tomocy/tapioca/cmd/tapioca/client/view/ascii"
	"github.com/tomocy/tapioca/domain"
)

func newCLI(cnf config) *CLI {
	return &CLI{
		cnf: cnf,
	}
}

type CLI struct {
	cnf config
}

func (c *CLI) Run() error {
	return c.summarizeCommitsOfToday()
}

func (c *CLI) summarizeCommitsOfToday() error {
	report := reportFunc("summarize commits of today")

	uc := newCommitUsecase()
	s, err := uc.SummarizeCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name)
	if err != nil {
		return report(err)
	}

	c.showSummary(s)

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

func (*CLI) showSummary(s *domain.Summary) {
	fmt.Println(ascii.ColorizedSummary(*s))
}

func (*CLI) showCommits(cs domain.Commits) {
	for _, c := range cs {
		fmt.Println(c)
	}
}
