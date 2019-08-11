package client

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/infra"
)

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

func parseConfig() (*config, error) {
	var r string
	flag.StringVar(&r, "r", "", "name of owner/repo")
	flag.Parse()

	splited := strings.Split(r, "/")
	if len(splited) != 2 {
		return nil, errors.New("invalid format of repo: the format of the name should be owner/repo")
	}

	return &config{
		repo: repo{
			owner: splited[0],
			name:  splited[1],
		},
	}, nil
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
