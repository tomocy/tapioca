package client

import (
	"net/url"
	"os"

	"github.com/pkg/browser"
	"github.com/tomocy/tapioca/domain"
)

func newCLI(printer printer) *CLI {
	return &CLI{
		printer: printer,
	}
}

type CLI struct {
	printer printer
}

func (c *CLI) ShowSummary(s domain.Summary) {
	c.printer.PrintSummary(os.Stdout, s)
}

type Twitter struct{}

func (t *Twitter) ShowSummary(s domain.Summary) {
	parsed, _ := url.Parse("https://twitter.com/intent/tweet")
	parsed.RawQuery = url.Values{
		"text":     []string{s.String()},
		"hashtags": []string{"tapioca"},
	}.Encode()

	browser.OpenURL(parsed.String())
}
