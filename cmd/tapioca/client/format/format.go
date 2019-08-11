package format

import (
	"fmt"
	"io"

	"github.com/tomocy/tapioca/cmd/tapioca/client/format/ascii"
	"github.com/tomocy/tapioca/domain"
)

type Text struct{}

func (t *Text) PrintSummary(w io.Writer, s domain.Summary) {
	fmt.Fprintln(w, s)
}

type Color struct{}

func (c *Color) PrintSummary(w io.Writer, s domain.Summary) {
	fmt.Fprintln(w, ascii.ColorizedSummary(s))
}
