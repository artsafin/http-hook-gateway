package structs

import (
	"http-hook-gateway/internal/parser"
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
	rawRequests    parser.RequestProvider
}

func (d *HookDef) MatchesAcceptUrl(acceptUrl string) bool {
	matches, err := regexp.Match(d.AcceptUrlRegex, []byte(acceptUrl))

	return err != nil && matches
}

func (d *HookDef) ParseRequest(req *http.Request) (parser.RequestProvider, error) {
	if d.rawRequests == nil {
		if parseErr := d.parseFile(); parseErr != nil {
			return nil, parseErr
		}
	}

	return nil, nil
}

func (d *HookDef) parseFile() error {
	file, openErr := os.Open(d.RequestFile)
	if openErr != nil {
		return openErr
	}
	defer file.Close()

	var err error
	d.rawRequests, err = parser.ParseFromReader(file)

	return err
}
