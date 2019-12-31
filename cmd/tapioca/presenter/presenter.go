package presenter

import (
	"io"
	"os"

	"github.com/tomocy/tapioca/domain"
)

type Stdout struct {
	Printer Printer
}

func (p *Stdout) PresentSummaries(ss ...*domain.Summary) {
	if len(ss) < 1 {
		return
	}
	if len(ss) == 1 {
		p.Printer.PrintSummary(os.Stdout, ss[0])
	}

	p.Printer.PrintSummaries(os.Stdout, ss)
}

type Printer interface {
	PrintSummaries(io.Writer, []*domain.Summary)
	PrintSummary(io.Writer, *domain.Summary)
}
