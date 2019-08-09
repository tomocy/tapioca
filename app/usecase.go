package app

import "github.com/tomocy/tapioca/domain"

type CommitUsecase struct {
	repo domain.CommitRepo
}

func (u *CommitUsecase) FetchCommits(owner, repo string) (domain.Commits, error) {
	return u.repo.FetchCommits(owner, repo)
}
