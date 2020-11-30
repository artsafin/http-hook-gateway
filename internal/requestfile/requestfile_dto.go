package requestfile

import (
	"io"
)

type dto struct {
	method  string
	path    string
	headers map[string]string
	body    io.Reader
}

func NewDto(method string, path string, headers map[string]string, body io.Reader) dto {
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
	return r.body
}
