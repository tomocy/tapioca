package client

import (
	"fmt"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/infra"
)

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

type config struct {
	repo repo
}

type repo struct {
	owner, name string
}

func reportFunc(did string) func(err error) error {
	return func(err error) error {
		return report(did, err)
	}
}

func report(did string, err error) error {
	return fmt.Errorf("failed to %s: %s", did, err)
}
