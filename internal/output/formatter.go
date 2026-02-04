package output

import (
	"io"

	"github.com/tehnerd/vape/internal/config"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

type Formatter interface {
	Format(data interface{}, writer io.Writer) error
}

func GetFormatter() Formatter {
	format := config.GetOutputFormat()
	switch Format(format) {
	case FormatJSON:
		return &JSONFormatter{}
	default:
		return &TableFormatter{}
	}
}

func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	default:
		return &TableFormatter{}
	}
}
