package domain

import "fmt"

type Commits []*Commit

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
