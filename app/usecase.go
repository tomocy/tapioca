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
	return u.fetchAndSummarizeCommits(owner, repo, domain.Params{
		Since: today(),
	})
}

func (u *CommitUsecase) SummarizeAuthorCommitsOfToday(owner, repo, author string) (*domain.Summary, error) {
	return u.fetchAndSummarizeCommits(owner, repo, domain.Params{
		Author: author,
		Since:  today(),
	})
}

func (u *CommitUsecase) fetchAndSummarizeCommits(owner, repo string, params domain.Params) (*domain.Summary, error) {
	today := today()
	s := &domain.Summary{
		Repo: &domain.Repo{
			Owner: owner,
			Name:  repo,
		},
		Date: today,
	}
	cs, err := u.repo.FetchCommits(owner, repo, params)
	if err != nil {
		return nil, err
	}
	s.Commits = cs
	s.Diff = cs.Diff()

	return s, nil
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
