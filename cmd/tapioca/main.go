package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tomocy/tapioca/app"
	"github.com/tomocy/tapioca/domain"
	"github.com/tomocy/tapioca/infra"
)

func main() {
	c := newClient()
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}

func newClient() client {
	conf := parseConfig()
	return &ofRepo{
		conf: conf,
		presenter: &stdout{
			printer: &text{
				colorized: conf.colorized,
			},
		},
	}
}

type client interface {
	Run() error
}

type ofRepo struct {
	conf      config
	presenter presenter
}

func (c *ofRepo) Run() error {
	s, err := c.summarize()
	if err != nil {
		return fmt.Errorf("failed to summarize: %s", err)
	}

	c.presenter.PresentSummary(*s)

	return nil
}

func (c *ofRepo) summarize() (*domain.Summary, error) {
	u := newCommitUsecase()
	return u.SummarizeCommits(c.conf.owner, c.conf.repo, domain.Params{
		Since: c.conf.since,
		Until: c.conf.until,
	})
}

func newCommitUsecase() *app.CommitUsecase {
	return app.NewCommitUsecase(infra.NewGitHub())
}

type help struct {
	err error
}

func (h help) Run() error {
	flag.Usage()
	return h.err
}

func parseConfig() config {
	owner, repo := flag.String("owner", "", "name of owner"), flag.String("repo", "", "name of repo")

	var since, until parseableTime
	flag.Var(&since, "since", "the day since which commits are summarized")
	flag.Var(&until, "until", "the day until which commits are summarized")

	colorized := flag.Bool("color", false, "colorize the output if true")

	flag.Parse()

	return config{
		owner: *owner, repo: *repo,
		since: time.Time(since), until: time.Time(until),
		colorized: *colorized,
	}
}

type parseableTime time.Time

func (t *parseableTime) Set(raw string) error {
	parsed, err := time.Parse("2006/01/02", raw)
	if err != nil {
		return err
	}

	*t = parseableTime(parsed)

	return nil
}

func (t parseableTime) String() string {
	return time.Time(t).Format("2006/01/02")
}

type config struct {
	owner        string
	repo         string
	since, until time.Time
	colorized    bool
}

type presenter interface {
	PresentSummary(domain.Summary)
}

type stdout struct {
	printer printer
}

func (p *stdout) PresentSummary(s domain.Summary) {
	p.printer.PrintSummary(os.Stdout, s)
}

type printer interface {
	PrintSummary(io.Writer, domain.Summary)
}

type text struct {
	colorized bool
}

func (p *text) PrintSummary(w io.Writer, s domain.Summary) {
	if p.colorized {
		fmt.Fprintln(w, colorizedSummary(s))
		return
	}

	fmt.Fprintln(w, s)
}

type colorizedSummary domain.Summary

func (s colorizedSummary) String() string {
	white := color.New(color.FgWhite).FprintfFunc()

	var b strings.Builder
	white(
		&b, "summary of commits to %s in %s\n%s",
		s.Repo, s.Since.Format("2006/01/02"), colorizedDiff(*s.Diff),
	)

	return b.String()
}

type colorizedDiff domain.Diff

func (d colorizedDiff) String() string {
	green, red := color.New(color.FgGreen).FprintfFunc(), color.New(color.FgRed).FprintfFunc()

	var b strings.Builder
	fmt.Fprintf(&b, "%d changes: ", d.Changes)
	green(&b, "%d adds", d.Adds)
	fmt.Fprint(&b, ", ")
	red(&b, "%d dels", d.Dels)

	return b.String()
}
