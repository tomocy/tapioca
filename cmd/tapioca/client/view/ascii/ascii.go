package ascii

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tomocy/tapioca/domain"
)

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
