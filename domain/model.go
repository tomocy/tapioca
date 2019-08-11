package domain

import (
	"fmt"
	"time"
)

type Summary struct {
	Repo    *Repo
	Commits []*Commit
	Diff    *Diff
	Date    time.Time
}

func (s Summary) String() string {
	return fmt.Sprintf(
		"summary of commits to %s in %s\n%s",
		s.Repo, s.Date.Format("2006/01/02"), s.Diff,
	)
}

type Repo struct {
	Owner, Name string
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type Commits []*Commit

func (cs Commits) Diff() *Diff {
	diff := new(Diff)
	diff.marge(cs.diffs()...)

	return diff
}

func (cs Commits) diffs() []*Diff {
	ds := make([]*Diff, len(cs))
	for i, c := range cs {
		ds[i] = c.Diff
	}

	return ds
}

type Commit struct {
	ID   string
	Diff *Diff
}

func (c Commit) String() string {
	return fmt.Sprintf("%s: %s", c.ID, c.Diff)
}

type Diff struct {
	Changes, Adds, Dels int
}

func (d Diff) String() string {
	return fmt.Sprintf("%d changes: %d adds, %d dels", d.Changes, d.Adds, d.Dels)
}

func (d *Diff) marge(ts ...*Diff) {
	for _, t := range ts {
		d.Changes += t.Changes
		d.Adds += t.Adds
		d.Dels += t.Dels
	}
}

type Params struct {
	Since time.Time
}
