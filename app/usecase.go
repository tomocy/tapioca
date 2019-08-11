package app

import (
	"time"

	"github.com/tomocy/tapioca/domain"
)

func NewCommitUsecase(repo domain.CommitRepo) *CommitUsecase {
	return &CommitUsecase{
		repo: repo,
	}
}

type CommitUsecase struct {
	repo domain.CommitRepo
}

func (u *CommitUsecase) SummarizeCommitsOfToday(owner, repo string) (*domain.Summary, error) {
	today := today()
	sum := &domain.Summary{
		Repo: &domain.Repo{
			Owner: owner,
			Name:  repo,
		},
		Date: today,
	}
	cs, err := u.repo.FetchCommits(owner, repo, &domain.Params{
		Since: today,
	})
	if err != nil {
		return nil, err
	}
	sum.Commits = cs
	sum.Diff = cs.Diff()

	return sum, nil
}

func (u *CommitUsecase) FetchCommitsOfToday(owner, repo string) (domain.Commits, error) {
	return u.repo.FetchCommits(owner, repo, nil)
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
