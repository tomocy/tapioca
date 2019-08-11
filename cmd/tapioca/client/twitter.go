package client

import (
	"net/url"

	"github.com/pkg/browser"
	"github.com/tomocy/tapioca/domain"
)

type Twitter struct {
	cnf config
}

func (t *Twitter) ShowSummary(s domain.Summary) {
	parsed, _ := url.Parse("https://twitter.com/intent/tweet")
	parsed.RawQuery = url.Values{
		"text": []string{s.String()},
	}.Encode()

	browser.OpenURL(parsed.String())
}
