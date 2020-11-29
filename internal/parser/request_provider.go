package parser

import (
	"io"
)

type RequestProvider interface {
	Method() string
	Path(defaultPath string) string
	Headers() map[string]string
	Body() io.Reader
}
