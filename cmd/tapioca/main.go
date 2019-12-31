package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tomocy/tapioca/app"
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

type client interface {
	Run() error
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

func parseConfig() config {
	owner, repo := flag.String("owner", "", "name of owner"), flag.String("repo", "", "name of repo")

	var since, until parseableTime
	flag.Var(&since, "since", "the day since which commits are summarized")
	flag.Var(&until, "until", "the day until which commits are summarized")

	colorized := flag.Bool("color", false, "colorize the output if true")

	flag.Parse()

	return config{
		owner: *owner, repo: *repo,
		since: time.Time(since), until: time.Time(until),
		colorized: *colorized,
	}
}

type parseableTime time.Time

func (t *parseableTime) Set(raw string) error {
	parsed, err := time.Parse("2006/01/02", raw)
	if err != nil {
		return err
	}

	*t = parseableTime(parsed)

	return nil
}

func (t parseableTime) String() string {
	return time.Time(t).Format("2006/01/02")
}

type config struct {
	owner        string
	repo         string
	since, until time.Time
	colorized    bool
}

type presenter interface {
	PresentSummary(domain.Summary)
}
