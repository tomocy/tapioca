package domain

type CommitRepo interface {
	FetchCommits(owner, repo string, params *Params) (Commits, error)
}
