package client

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func New() Client {
	cnf, err := parseConfig()
	if err != nil {
		return &Help{
			err: err,
		}
	}

	var client Client
	switch cnf.mode {
	case modeCLI:
		client = newCLI(*cnf)
	default:
		client = new(Help)
	}

	return client
}

type Client interface {
	Run() error
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

func parseConfig() (*config, error) {
	m := flag.String("m", string(modeCLI), "name of mode")
	r := flag.String("r", "", "name of owner/repo")
	flag.Parse()

	cnf := new(config)
	if err := cnf.parseRepo(*r); err != nil {
		return nil, err
	}
	cnf.mode = mode(*m)

	return cnf, nil
}

type config struct {
	mode mode
	repo repo
}

func (c *config) parseRepo(r string) error {
	splited := strings.Split(r, "/")
	if len(splited) != 2 {
		return errors.New("invalid format of repo: the format of the name should be owner/repo")
	}

	c.repo = repo{
		owner: splited[0],
		name:  splited[1],
	}

	return nil
}

const (
	modeCLI mode = "cli"
)

type mode string

type format string

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

type Help struct {
	err error
}

func (h *Help) Run() error {
	flag.Usage()
	return h.err
}

type printer interface {
	PrintSummary(io.Writer, domain.Summary)
}
