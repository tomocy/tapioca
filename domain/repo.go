package domain

import "time"

type CommitRepo interface {
	FetchCommitsSinceDate(owner, repo string, date time.Time) (Commits, error)
}
