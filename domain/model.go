package domain

import (
	"sort"
	"time"
)

type Summary struct {
	Repo         *Repo
	Authors      []string
	Commits      []*Commit
	Diff         *Diff
	Since, Until time.Time
}

type Repo struct {
	Owner, Name string
}

type Commits []*Commit

func (cs Commits) Authors() []string {
	am := make(map[string]bool)
	for _, c := range cs {
		am[c.Author] = true
	}
	as := make([]string, len(am))
	var i int
	for a := range am {
		as[i] = a
		i++
	}
	sort.Slice(as, func(i, j int) bool {
		return as[i] < as[j]
	})

	return as
}

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
	ID        string
	Author    string
	Diff      *Diff
	CreatedAt time.Time
}

type Diff struct {
	Changes, Adds, Dels int
}

func (d *Diff) marge(ts ...*Diff) {
	for _, t := range ts {
		d.Changes += t.Changes
		d.Adds += t.Adds
		d.Dels += t.Dels
	}
}

type Params struct {
	Author       string
	Since, Until time.Time
}
