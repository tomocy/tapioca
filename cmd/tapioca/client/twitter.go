package client

import (
	"net/url"

	"github.com/pkg/browser"
	"github.com/tomocy/tapioca/domain"
)

func newTwitter(cnf config) *Twitter {
	return &Twitter{
		cnf: cnf,
	}
}

type Twitter struct {
	cnf config
}

func (t *Twitter) Run() error {
	if t.cnf.author != "" {
		return t.summarizeAuthorCommitsOfToday()
	}

	return t.summarizeCommitsOfToday()
}

func (t *Twitter) summarizeAuthorCommitsOfToday() error {
	s, err := summarizeAuthorCommitsOfToday(t.cnf.repo.owner, t.cnf.repo.name, t.cnf.author)
	if err != nil {
		return err
	}

	t.showSummary(*s)

	return nil
}

func (t *Twitter) summarizeCommitsOfToday() error {
	s, err := summarizeCommitsOfToday(t.cnf.repo.owner, t.cnf.repo.name)
	if err != nil {
		return err
	}

	t.showSummary(*s)

	return nil
}

func (t *Twitter) ShowSummary(s domain.Summary) {
	parsed, _ := url.Parse("https://twitter.com/intent/tweet")
	parsed.RawQuery = url.Values{
		"text": []string{s.String()},
	}.Encode()

	browser.OpenURL(parsed.String())
}

func (t *Twitter) showSummary(s domain.Summary) {
	parsed, _ := url.Parse("https://twitter.com/intent/tweet")
	parsed.RawQuery = url.Values{
		"text": []string{s.String()},
	}.Encode()

	browser.OpenURL(parsed.String())
}
