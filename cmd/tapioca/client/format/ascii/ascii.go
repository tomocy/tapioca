package ascii

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tomocy/tapioca/domain"
)

type ColorizedSummary domain.Summary

func (s ColorizedSummary) String() string {
	white := color.New(color.FgWhite).FprintfFunc()
	var b strings.Builder
	white(
		&b, "summary of commits to %s in %s\n%s",
		s.Repo, s.Since.Format("2006/01/02"), ColorizedDiff(*s.Diff),
	)

	return b.String()
}

type ColorizedDiff domain.Diff

func (d ColorizedDiff) String() string {
	green, red := color.New(color.FgGreen).FprintfFunc(), color.New(color.FgRed).FprintfFunc()
	var b strings.Builder
	fmt.Fprintf(&b, "%d changes: ", d.Changes)
	green(&b, "%d adds", d.Adds)
	fmt.Fprint(&b, ", ")
	red(&b, "%d dels", d.Dels)

	return b.String()
}
