package application

import (
	"encoding/json"
	"errors"
	"http-hook-gateway/internal/requestfile"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var ErrSkipRequest = errors.New("request is skipped")

func interpolateRequestfile(src requestfile.RequestFile, data *RequestSummary) (requestfile.RequestFile, error) {
	method, methodErr := renderString("method", src.Method(), data)
	if methodErr != nil {
		return nil, methodErr
	}

	if len(method) == 0 {
		return nil, ErrSkipRequest
	}

	path, pathErr := renderString("path", src.Path("/"), data)
	if pathErr != nil {
		return nil, pathErr
	}

	headers, headersErr := renderMap("headers", src.Headers(), data)
	if headersErr != nil {
		return nil, headersErr
	}

	body, bodyErr := renderReader("body", src.Body(), data)
	if bodyErr != nil {
		return nil, bodyErr
	}

	dto := requestfile.NewDto(method, path, headers, body)

	return dto, nil
}

func newTpl(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"json": func(data interface{}) string {
			bytes, jsonErr := json.Marshal(data)
			if jsonErr != nil {
				return ""
			}
			return string(bytes)
		},
		"query": func(data *RequestSummary, name string) string {
			return data.Query.Get(name)
		},
		"header": func(data *RequestSummary, name string) string {
			return data.Headers.Get(name)
		},
		"headervalues": func(data *RequestSummary, name string) string {
			return strings.Join(data.Headers.Values(name), "; ")
		},
		"env": func(key string) string {
			return os.Getenv(key)
		},
	})
}

func renderString(name, tpl string, summary *RequestSummary) (string, error) {
	var b strings.Builder
	err := template.Must(newTpl(name).Parse(tpl)).Execute(&b, summary)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func renderMap(name string, m map[string]string, r *RequestSummary) (map[string]string, error) {
	res := make(map[string]string)
	for k, v := range m {
		var err error
		res[k], err = renderString(name+k, v, r)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func renderReader(name string, reader io.Reader, r *RequestSummary) (string, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	res, err := renderString(name, string(bytes), r)
	if err != nil {
		return "", err
	}

	return res, nil
}
