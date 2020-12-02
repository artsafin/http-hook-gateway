package requestfile

import (
	"io"
	"strings"
)

type dto struct {
	method  string
	path    string
	headers map[string]string
	body    string
}

func NewDto(method string, path string, headers map[string]string, body string) dto {
	return dto{method: method, path: path, headers: headers, body: body}
}

func (r dto) Method() string {
	return r.method
}

func (r dto) Path(defaultPath string) string {
	if r.path == "" {
		return defaultPath
	}
	return r.path
}

func (r dto) Headers() map[string]string {
	return r.headers
}

func (r dto) Body() io.Reader {
	return strings.NewReader(r.body)
}
