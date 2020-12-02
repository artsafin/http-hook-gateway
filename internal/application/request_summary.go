package application

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type RequestSummary struct {
	Method       string
	RemoteAddr   string
	Headers      http.Header
	Body         interface{}
	Scheme       string
	User         string
	UserPassword string
	Hostname     string
	Path         string
	Fragment     string
	Query        url.Values
}

func NewSummaryFromHttp(req *http.Request) (*RequestSummary, error) {
	summary := RequestSummary{
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		Headers:    req.Header,
		Body:       make(map[string]string),
		Scheme:     req.URL.Scheme,
		User:       req.URL.User.Username(),
		Hostname:   req.Host,
		Path:       req.URL.Path,
		Fragment:   req.URL.Fragment,
		Query:      req.URL.Query(),
	}
	if userPassword, passwordSet := req.URL.User.Password(); passwordSet {
		summary.UserPassword = userPassword
	}

	ct := req.Header.Get("Content-Type")
	if ct == "application/json" {
		var jsonErr error
		summary.Body, jsonErr = getJsonBodyData(req.Body)
		if jsonErr != nil {
			return nil, jsonErr
		}
		defer req.Body.Close()
	} else {
		formErr := req.ParseMultipartForm(1 * 1024 * 1024 * 1024)

		if formErr != nil && formErr != http.ErrNotMultipart {
			return nil, formErr
		}

		body := make(map[string]string)

		for k, v := range req.PostForm {
			body[k] = strings.Join(v, "; ")
		}

		summary.Body = body
	}

	return &summary, nil
}

func getJsonBodyData(body io.Reader) (interface{}, error) {
	reqBody, reqReadErr := ioutil.ReadAll(body)
	if reqReadErr != nil {
		return nil, reqReadErr
	}

	data := make(map[string]interface{})
	jsonErr := json.Unmarshal(reqBody, &data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return data, nil
}
