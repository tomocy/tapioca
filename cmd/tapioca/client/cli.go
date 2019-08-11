package client

import (
	"flag"
	"fmt"

	"github.com/tomocy/tapioca/domain"
)

type CLI struct{}

func (c *CLI) Run() error {
	return c.summarizeCommitsOfToday()
}

func (c *CLI) summarizeCommitsOfToday() error {
	report := reportFunc("summarize commits of today")
	cnf, err := parseConfig()
	if err != nil {
		flag.Usage()
		return report(err)
	}

	uc := newCommitUsecase()
	s, err := uc.SummarizeCommitsOfToday(cnf.repo.owner, cnf.repo.name)
	if err != nil {
		return report(err)
	}

	c.showSummary(s)

	return nil
}

func (c *CLI) fetchCommits() error {
	report := reportFunc("fetch commits")
	cnf, err := parseConfig()
	if err != nil {
		flag.Usage()
		return report(err)
	}

	uc := newCommitUsecase()
	cs, err := uc.FetchCommitsOfToday(cnf.repo.owner, cnf.repo.name)
	if err != nil {
		return report(err)
	}

	c.showCommits(cs)

	return nil
}

func (*CLI) showSummary(s *domain.Summary) {
	fmt.Println(s)
}

func (*CLI) showCommits(cs domain.Commits) {
	for _, c := range cs {
		fmt.Println(c)
	}
}
