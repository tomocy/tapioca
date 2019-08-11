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

func (u *CommitUsecase) FetchCommitsOfToday(owner, repo string) (domain.Commits, error) {
	return u.repo.FetchCommitsSinceDate(owner, repo, today())
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
