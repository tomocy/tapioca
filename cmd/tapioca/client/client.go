package client

type config struct {
	repo repo
}

type repo struct {
	owner, name string
}
