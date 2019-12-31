package printer

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tomocy/tapioca/domain"
)

type InText struct {
	Colorized bool
}

func (p *InText) PrintSummary(w io.Writer, s domain.Summary) {
	target := colorizedSummary{
		Summary: s,
		white:   fmt.Fprintf,
		green:   fmt.Fprintf,
		red:     fmt.Fprintf,
	}
	if p.Colorized {
		target.white = color.New(color.FgWhite).Fprintf
		target.green = color.New(color.FgGreen).Fprintf
		target.red = color.New(color.FgRed).Fprintf
	}

	fmt.Fprintln(w, target)
}

type colorizedSummary struct {
	domain.Summary
	white func(io.Writer, string, ...interface{}) (int, error)
	green func(io.Writer, string, ...interface{}) (int, error)
	red   func(io.Writer, string, ...interface{}) (int, error)
}

func (s colorizedSummary) String() string {
	su := sinceUntil{
		since: s.Since,
		until: s.Until,
	}

	var b strings.Builder
	s.white(&b, "summary of commits to %s/%s %s\n", s.Repo.Owner, s.Repo.Name, su.Format("2006/01/02"))
	s.white(&b, "%d commits\n", len(s.Commits))
	s.white(&b, "%d changes: ", s.Diff.Changes)
	s.green(&b, "%d adds", s.Diff.Adds)
	s.white(&b, ", ")
	s.red(&b, "%d dels", s.Diff.Dels)

	return b.String()
}

type sinceUntil struct {
	since, until time.Time
}

func (su sinceUntil) Format(format string) string {
	if !su.since.IsZero() && !su.until.IsZero() {
		return fmt.Sprintf("from %s to %s", su.since.Format(format), su.until.Format(format))
	}

	if !su.since.IsZero() {
		return fmt.Sprintf("since %s", su.since.Format(format))
	}

	return fmt.Sprintf("until %s", su.until.Format(format))
}
