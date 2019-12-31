package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/cmd/tapioca/client"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func main() {
	c := client.New()
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}

type ofRepo struct {
	conf      config
	presenter presenter
}

func (c *ofRepo) Run() error {
	s, err := c.summarize()
	if err != nil {
		return fmt.Errorf("failed to summarize: %s", err)
	}

	c.presenter.PresentSummary(*s)

	return nil
}

func (c *ofRepo) summarize() (*domain.Summary, error) {
	u := newCommitUsecase()
	return u.SummarizeCommits(c.conf.owner, c.conf.repo, domain.Params{
		Since: c.conf.since,
		Until: c.conf.until,
	})
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
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
