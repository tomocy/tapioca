package domain

type RepoRepo interface {
	FetchRepos(owner string) ([]*Repo, error)
}

type CommitRepo interface {
	FetchCommits(owner, repo string, params Params) (Commits, error)
}
