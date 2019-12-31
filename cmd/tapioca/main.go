package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientPkg "github.com/tomocy/tapioca/cmd/tapioca/client"
	"github.com/tomocy/tapioca/cmd/tapioca/presenter"
	"github.com/tomocy/tapioca/cmd/tapioca/printer"
)

func main() {
	conf := parseConfig()
	c := newClient(conf)

	ctx := contextWithSignals(context.Background(), syscall.SIGINT)

	if err := c.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
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

func newClient(conf config) client {
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

func contextWithSignals(ctx context.Context, sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		defer close(sigCh)
		defer cancel()

		sig := <-sigCh
		fmt.Println(sig)
		return
	}()

	return ctx
}

type config struct {
	owner        string
	repo         string
	author       string
	since, until time.Time
	colorized    bool
}

type client interface {
	Run(context.Context) error
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
