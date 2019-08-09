package domain

type CommitRepo interface {
	FetchCommits(owner, repo string) (Commits, error)
}
