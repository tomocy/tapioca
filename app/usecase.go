package app

import "github.com/tomocy/tapioca/domain"

func NewCommitUsecase(repo domain.CommitRepo) *CommitUsecase {
	return &CommitUsecase{
		repo: repo,
	}
}

type CommitUsecase struct {
	repo domain.CommitRepo
}

func (u *CommitUsecase) FetchCommits(owner, repo string) (domain.Commits, error) {
	return u.repo.FetchCommits(owner, repo)
}
