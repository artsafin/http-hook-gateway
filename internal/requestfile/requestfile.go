package requestfile

import (
	"io"
)

type RequestFile interface {
	Method() string
	Path(defaultPath string) string
	Headers() map[string]string
	Body() io.Reader
}
