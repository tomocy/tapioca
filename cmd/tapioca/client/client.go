package client

import (
	"fmt"
	"time"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func NewOfRepos(owner, author string, since, until time.Time, presenter Presenter) *OfRepos {
	return &OfRepos{
		owner: owner, author: author,
		since: since, until: until,
		presenter: presenter,
	}
}

type OfRepos struct {
	owner, author string
	since, until  time.Time
	presenter     Presenter
}

func (c *OfRepos) Run() error {
	ss, err := c.summarize()
	if err != nil {
		return err
	}

	for _, s := range ss {
		c.presenter.PresentSummary(*s)
	}

	return nil
}

func (c *OfRepos) summarize() ([]*domain.Summary, error) {
	u := newRepoUsecase()
	return u.SummarizeCommits(c.owner, domain.Params{
		Author: c.author,
		Since:  c.since,
		Until:  c.until,
	})
}

func NewOfRepo(owner, repo, author string, since, until time.Time, presenter Presenter) *OfRepo {
	return &OfRepo{
		owner: owner, repo: repo, author: author,
		since: since, until: until,
		presenter: presenter,
	}
}

type OfRepo struct {
	owner, repo, author string
	since, until        time.Time
	presenter           Presenter
}

func (c *OfRepo) Run() error {
	s, err := c.summarize()
	if err != nil {
		return fmt.Errorf("failed to summarize: %s", err)
	}

	c.presenter.PresentSummary(*s)

	return nil
}

func (c *OfRepo) summarize() (*domain.Summary, error) {
	u := newCommitUsecase()
	return u.SummarizeCommits(c.owner, c.repo, domain.Params{
		Author: c.author,
		Since:  c.since,
		Until:  c.until,
	})
}

func newRepoUsecase() *app.RepoUsecase {
	return app.NewRepoUsecase(githubRepo, newCommitUsecase())
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(githubRepo)
}

var githubRepo = infra.NewGitHub()

type Presenter interface {
	PresentSummary(domain.Summary)
}
