package client

import (
	"fmt"
	"time"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func NewOfRepo(owner, repo string, since, until time.Time, presenter Presenter) *OfRepo {
	return &OfRepo{
		owner: owner, repo: repo,
		since: since, until: until,
		presenter: presenter,
	}
}

type OfRepo struct {
	owner, repo  string
	since, until time.Time
	presenter    Presenter
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
		Since: c.since,
		Until: c.until,
	})
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

type Presenter interface {
	PresentSummary(domain.Summary)
}
