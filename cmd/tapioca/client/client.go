package client

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/tomocy/tapioca/app"
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
	u := newCommitUsecase()
	s, err := u.SummarizeCommits(c.cnf.repo.owner, c.cnf.repo.name, domain.Params{
		Author: c.cnf.author,
		Since:  c.cnf.since,
		Until:  c.cnf.until,
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func parseConfig() (*config, error) {
	m := flag.String("m", modeCLI, "name of mode")
	f := flag.String("f", formatText, "name of format")
	d := flag.String("d", dayToday, "day of commits to be summarized")
	r := flag.String("r", "", "name of owner/repo")
	a := flag.String("a", "", "name of author")
	flag.Parse()

	cnf := &config{
		mode: *m, format: *f,
		author: *a,
	}

	if *d == dayToday {
		cnf.since = today()
	} else if *d == dayYesterday {
		cnf.since, cnf.until = yesterday(), today()
	}

	if err := cnf.parseRepo(*r); err != nil {
		return nil, err
	}

	return cnf, nil
}

func yesterday() time.Time {
	return today().Add(-24 * time.Hour)
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

type config struct {
	mode         string
	format       string
	since, until time.Time
	repo         repo
	author       string
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
	var p printer = new(Text)
	switch fmt {
	case formatColor:
		p = new(Color)
	}

	return p
}

type printer interface {
	PrintSummary(io.Writer, domain.Summary)
}

const (
	formatText  = "text"
	formatColor = "color"

	dayToday     = "today"
	dayYesterday = "yesterday"
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
