package application

import (
	"http-hook-gateway/internal/requestfile"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

func transformRequestfile(src requestfile.RequestFile, summary RequestSummary) (requestfile.RequestFile, error) {
	method, methodErr := renderString("method", src.Method(), &summary)
	if methodErr != nil {
		return nil, methodErr
	}

	path, pathErr := renderString("path", src.Path("/"), &summary)
	if pathErr != nil {
		return nil, pathErr
	}

	headers, headersErr := renderMap("headers", src.Headers(), &summary)
	if headersErr != nil {
		return nil, headersErr
	}

	body, bodyErr := renderReader("body", src.Body(), &summary)
	if bodyErr != nil {
		return nil, bodyErr
	}

	dto := requestfile.NewDto(method, path, headers, body)

	return dto, nil
}

func renderString(name, tpl string, summary *RequestSummary) (string, error) {
	var b strings.Builder
	err := template.Must(template.New(name).Parse(tpl)).Execute(&b, summary)
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

func renderReader(name string, reader io.Reader, r *RequestSummary) (io.Reader, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res, err := renderString(name, string(bytes), r)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(res), nil
}
