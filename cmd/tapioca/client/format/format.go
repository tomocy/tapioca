package format

import (
	"fmt"
	"io"

	"github.com/tomocy/tapioca/domain"
)

type Text struct{}

func (t *Text) PrintSummary(w io.Writer, s domain.Summary) {
	fmt.Fprintln(w, s)
}
