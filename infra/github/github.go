package github

import (
	"github.com/tomocy/tapioca/domain"
)

type Commits []*Commit

func (cs Commits) Adapt() []*domain.Commit {
	adapteds := make([]*domain.Commit, len(cs))
	for i, c := range cs {
		adapteds[i] = c.Adapt()
	}

	return adapteds
}

type Commit struct {
	SHA   string  `json:"sha"`
	Files []*Diff `json:"files"`
}

func (c *Commit) Adapt() *domain.Commit {
	return &domain.Commit{
		ID:   c.SHA,
		Diff: c.adaptDiff(),
	}
}

func (c *Commit) adaptDiff() *domain.Diff {
	if len(c.Files) <= 0 {
		return nil
	}

	return c.Files[0].Adapt()
}

type Diff struct {
	Changes   int `json:"changes"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

func (d *Diff) Adapt() *domain.Diff {
	return &domain.Diff{
		Changes: d.Changes,
		Adds:    d.Additions,
		Dels:    d.Deletions,
	}
}
