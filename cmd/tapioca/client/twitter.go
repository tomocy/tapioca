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
	return t.summarizeCommitsOfToday()
}

func (t *Twitter) summarizeCommitsOfToday() error {
	report := reportFunc("summarize commits of today")

	uc := newCommitUsecase()
	s, err := uc.SummarizeCommitsOfToday(t.cnf.repo.owner, t.cnf.repo.name)
	if err != nil {
		return report(err)
	}

	t.showSummary(*s)

	return nil
}

func (t *Twitter) showSummary(s domain.Summary) {
	parsed, _ := url.Parse("https://twitter.com/intent/tweet")
	parsed.RawQuery = url.Values{
		"text": []string{s.String()},
	}.Encode()

	browser.OpenURL(parsed.String())
}
