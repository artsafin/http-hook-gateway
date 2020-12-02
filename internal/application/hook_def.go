package application

import (
	"http-hook-gateway/internal/requestfile"
	"net/http"
	"os"
	"regexp"
)

type HookDef struct {
	Name           string
	AcceptUrlRegex string
	ProxyHost      string
	ProxyPath      string
	RequestFile    string
	rawRequest     requestfile.RequestFile
}

func (d *HookDef) MatchesAcceptUrl(acceptUrl string) bool {
	matches, err := regexp.Match(d.AcceptUrlRegex, []byte(acceptUrl))

	return err == nil && matches
}

func (d *HookDef) ParseRequest(req *http.Request) (requestfile.RequestFile, error) {
	if d.rawRequest == nil {
		if parseErr := d.parseFile(); parseErr != nil {
			return nil, parseErr
		}
	}

	summary, summaryErr := NewSummaryFromHttp(req)
	if summaryErr != nil {
		return nil, summaryErr
	}

	return interpolateRequestfile(d.rawRequest, summary)
}

func (d *HookDef) parseFile() error {
	file, openErr := os.Open(d.RequestFile)
	if openErr != nil {
		return openErr
	}
	defer file.Close()

	var err error
	d.rawRequest, err = requestfile.ParseFromReader(file)

	return err
}
