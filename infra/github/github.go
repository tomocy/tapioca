package github

import "github.com/tomocy/tapioca/domain"

type Commits []*Commit

func (cs Commits) Adapt() []*domain.Commit {
	adapteds := make([]*domain.Commit, len(cs))
	for i, c := range cs {
		adapteds[i] = c.Adapt()
	}

	return adapteds
}

type Commit struct {
	SHA string `json:"sha"`
}

func (c *Commit) Adapt() *domain.Commit {
	return &domain.Commit{
		ID: c.SHA,
	}
}
