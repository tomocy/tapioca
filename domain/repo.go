package domain

import "context"

type RepoRepo interface {
	FetchRepos(context.Context, string) ([]*Repo, error)
}

type CommitRepo interface {
	FetchCommits(context.Context, string, string, Params) (Commits, error)
}
