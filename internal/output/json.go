package output

import (
	"encoding/json"
	"io"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(data interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
