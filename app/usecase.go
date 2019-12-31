package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tomocy/tapioca/domain"
)

func NewRepoUsecase(repo domain.RepoRepo, commit *CommitUsecase) *RepoUsecase {
	return &RepoUsecase{
		repo:   repo,
		commit: commit,
	}
}

type RepoUsecase struct {
	repo   domain.RepoRepo
	commit *CommitUsecase
}

func (u *RepoUsecase) SummarizeCommits(ctx context.Context, owner string, params domain.Params) ([]*domain.Summary, error) {
	repos, err := u.repo.FetchRepos(ctx, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %s", err)
	}

	ch := make(chan *domain.Summary)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		for _, repo := range repos {
			wg.Add(1)
			go func(repo *domain.Repo) {
				defer wg.Done()
				s, err := u.commit.SummarizeCommits(ctx, repo.Owner, repo.Name, params)
				if err != nil {
					return
				}

				ch <- s
			}(repo)
		}

		wg.Wait()
	}()

	ss := make([]*domain.Summary, 0, len(repos))
	for s := range ch {
		ss = append(ss, s)
	}

	return ss, nil
}

func NewCommitUsecase(repo domain.CommitRepo) *CommitUsecase {
	return &CommitUsecase{
		repo: repo,
	}
}

type CommitUsecase struct {
	repo domain.CommitRepo
}

func (u *CommitUsecase) SummarizeCommits(ctx context.Context, owner, repo string, params domain.Params) (*domain.Summary, error) {
	s := &domain.Summary{
		Repo: &domain.Repo{
			Owner: owner,
			Name:  repo,
		},
		Since: params.Since,
		Until: params.Until,
	}
	cs, err := u.repo.FetchCommits(ctx, owner, repo, params)
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
