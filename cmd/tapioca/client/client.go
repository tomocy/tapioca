package client

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/tomocy/tapioca/app"
	formatPkg "github.com/tomocy/tapioca/cmd/tapioca/client/format"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func New() Runner {
	cnf, err := parseConfig()
	if err != nil {
		return &Help{
			err: err,
		}
	}

	return &Client{
		cnf:       *cnf,
		presenter: newPresenter(*cnf),
	}
}

type Runner interface {
	Run() error
}

type Client struct {
	cnf       config
	presenter presenter
}

func (c *Client) Run() error {
	s, err := c.summarize()
	if err != nil {
		return err
	}

	c.presenter.ShowSummary(*s)

	return nil
}

func (c *Client) summarize() (*domain.Summary, error) {
	var s *domain.Summary
	var err error
	var report func(error) error
	if c.cnf.author != "" {
		s, err = summarizeAuthorCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name, c.cnf.author)
		report = reportFunc("summarize author commits of today")
	} else {
		s, err = summarizeCommitsOfToday(c.cnf.repo.owner, c.cnf.repo.name)
		report = reportFunc("summarize commits of today")
	}
	if err != nil {
		return nil, report(err)
	}

	return s, nil
}

func parseConfig() (*config, error) {
	m := flag.String("m", modeCLI, "name of mode")
	f := flag.String("f", formatText, "name of format")
	r := flag.String("r", "", "name of owner/repo")
	a := flag.String("a", "", "name of author")
	flag.Parse()

	cnf := &config{
		mode:   *m,
		format: *f,
		author: *a,
	}
	if err := cnf.parseRepo(*r); err != nil {
		return nil, err
	}

	return cnf, nil
}

type config struct {
	mode   string
	format string
	repo   repo
	author string
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

func newPresenter(cnf config) presenter {
	switch cnf.mode {
	case modeCLI:
		return newCLI(newPrinter(cnf.format))
	case modeTwitter:
		return new(Twitter)
	default:
		return new(Help)
	}
}

type presenter interface {
	ShowSummary(domain.Summary)
}

const (
	modeCLI     = "cli"
	modeTwitter = "twitter"
)

func newPrinter(fmt string) printer {
	var p printer = new(formatPkg.Text)
	switch fmt {
	case formatColor:
		p = new(formatPkg.Color)
	}

	return p
}

type printer interface {
	PrintSummary(io.Writer, domain.Summary)
}

const (
	formatText  = "text"
	formatColor = "color"
)

type repo struct {
	owner, name string
}

type Help struct {
	err error
}

func (h *Help) Run() error {
	flag.Usage()
	return h.err
}

func (h *Help) ShowSummary(domain.Summary) {
	flag.Usage()
}

func summarizeCommitsOfToday(owner, repo string) (*domain.Summary, error) {
	report := reportFunc("summarize commits of today")
	uc := newCommitUsecase()
	s, err := uc.SummarizeCommitsOfToday(owner, repo)
	if err != nil {
		return nil, report(err)
	}

	return s, nil
}

func summarizeAuthorCommitsOfToday(owner, repo, author string) (*domain.Summary, error) {
	report := reportFunc("summarize author commits of today")
	uc := newCommitUsecase()
	s, err := uc.SummarizeAuthorCommitsOfToday(owner, repo, author)
	if err != nil {
		return nil, report(err)
	}

	return s, nil
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

func reportFunc(did string) func(err error) error {
	return func(err error) error {
		return report(did, err)
	}
}

func report(did string, err error) error {
	return fmt.Errorf("failed to %s: %s", did, err)
}
