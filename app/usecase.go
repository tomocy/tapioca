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

func (u *CommitUsecase) SummarizeCommits(owner, repo string, params domain.Params) (*domain.Summary, error) {
	s := &domain.Summary{
		Repo: &domain.Repo{
			Owner: owner,
			Name:  repo,
		},
		Since: params.Since,
		Until: params.Until,
	}
	cs, err := u.repo.FetchCommits(owner, repo, params)
	if err != nil {
		return nil, err
	}
	s.Commits = cs
	s.Authors = cs.Authors()
	s.Diff = cs.Diff()

	return s, nil
}

func yesterday() time.Time {
	return today().Add(-24 * time.Hour)
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
