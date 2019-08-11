package client

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/tomocy/tapioca/domain"
)

type CLI struct{}

func (c *CLI) Run() error {
	return c.summarizeCommitsOfToday()
}

func (c *CLI) summarizeCommitsOfToday() error {
	report := reportFunc("summarize commits of today")
	cnf, err := c.parseConfig()
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
	cnf, err := c.parseConfig()
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

func (c *CLI) parseConfig() (*config, error) {
	var r string
	flag.StringVar(&r, "r", "", "name of owner/repo")
	flag.Parse()

	splited := strings.Split(r, "/")
	if len(splited) != 2 {
		return nil, errors.New("invalid format of repo: the format of the name should be owner/repo")
	}

	return &config{
		repo: repo{
			owner: splited[0],
			name:  splited[1],
		},
	}, nil
}

func (*CLI) showSummary(s *domain.Summary) {
	fmt.Println(s)
}

func (*CLI) showCommits(cs domain.Commits) {
	for _, c := range cs {
		fmt.Println(c)
	}
}
