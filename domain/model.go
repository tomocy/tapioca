package domain

import (
	"fmt"
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

func (s Summary) String() string {
	su := sinceUntil{
		since: s.Since,
		until: s.Until,
	}
	return fmt.Sprintf(
		"summary of commits to %s in %s\n%d commits\n%s",
		s.Repo, su.Format("2006/01/02"),
		len(s.Commits),
		s.Diff,
	)
}

type sinceUntil struct {
	since, until time.Time
}

func (su sinceUntil) Format(format string) string {
	if !su.since.IsZero() && !su.until.IsZero() {
		return fmt.Sprintf("from %s to %s", su.since.Format(format), su.until.Format(format))
	}

	if !su.since.IsZero() {
		return fmt.Sprintf("since %s", su.since.Format(format))
	}

	return fmt.Sprintf("until %s", su.until.Format(format))
}

type Repo struct {
	Owner, Name string
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
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

func (c Commit) String() string {
	return fmt.Sprintf("%s: %s: %s: %s", c.CreatedAt.Format("2006/01/02"), c.Author, c.ID, c.Diff)
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
	Author       string
	Since, Until time.Time
}
