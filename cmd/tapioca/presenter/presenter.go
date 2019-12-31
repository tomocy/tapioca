package presenter

import (
	"io"
	"os"

	"github.com/tomocy/tapioca/domain"
)

type Stdout struct {
	Printer Printer
}

func (p *Stdout) PresentSummary(s domain.Summary) {
	p.Printer.PrintSummary(os.Stdout, s)
}

type Printer interface {
	PrintSummary(io.Writer, domain.Summary)
}
