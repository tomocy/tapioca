package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/tomocy/tapioca/domain"
)

type InText struct {
	Colorized bool
}

func (p *InText) PrintSummary(w io.Writer, s domain.Summary) {
	if p.Colorized {
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
