package domain

type Commit struct {
	ID string
}

type Diff struct {
	Base, Head          *Commit
	Changes, Adds, Dels int
}
