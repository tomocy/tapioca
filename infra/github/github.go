package github

import (
	"time"

	"github.com/tomocy/tapioca/domain"
)

type Repos []*Repo

func (rs Repos) Adapt() []*domain.Repo {
	adapteds := make([]*domain.Repo, len(rs))
	for i, r := range rs {
		adapteds[i] = r.Adapt()
	}

	return adapteds
}

type Repo struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (r *Repo) Adapt() *domain.Repo {
	return &domain.Repo{
		Owner: r.Owner.Login,
		Name:  r.Name,
	}
}

type Commits []*Commit

func (cs Commits) Adapt() []*domain.Commit {
	adapteds := make([]*domain.Commit, len(cs))
	for i, c := range cs {
		adapteds[i] = c.Adapt()
	}

	return adapteds
}

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	}
	Files []*Diff `json:"files"`
}

func (c *Commit) Adapt() *domain.Commit {
	return &domain.Commit{
		ID:        c.SHA,
		Author:    c.Author.Login,
		Diff:      c.adaptDiff(),
		CreatedAt: c.Commit.Author.Date,
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

type Error struct {
	Message string `json:"message"`
}
