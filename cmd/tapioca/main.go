package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	clientPkg "github.com/tomocy/tapioca/cmd/tapioca/client"
	"github.com/tomocy/tapioca/cmd/tapioca/presenter"
	"github.com/tomocy/tapioca/cmd/tapioca/printer"
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

	if conf.owner != "" && conf.repo == "" {
		return clientPkg.NewOfRepos(
			conf.owner, conf.author,
			conf.since, conf.until,
			&presenter.Stdout{
				Printer: &printer.InText{
					Colorized: conf.colorized,
				},
			},
		)
	}

	return clientPkg.NewOfRepo(
		conf.owner, conf.repo, conf.author,
		conf.since, conf.until,
		&presenter.Stdout{
			Printer: &printer.InText{
				Colorized: conf.colorized,
			},
		},
	)
}

type client interface {
	Run() error
}

func parseConfig() config {
	owner, repo := flag.String("owner", "", "name of owner"), flag.String("repo", "", "name of repo")
	author := flag.String("author", "", "name or email address of author whose commits are summarized")

	var since, until parseableTime
	flag.Var(&since, "since", "the day since which commits are summarized")
	flag.Var(&until, "until", "the day until which commits are summarized")

	colorized := flag.Bool("color", false, "colorize the output if true")

	flag.Parse()

	return config{
		owner: *owner, repo: *repo,
		author: *author,
		since:  time.Time(since), until: time.Time(until),
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
	author       string
	since, until time.Time
	colorized    bool
}
