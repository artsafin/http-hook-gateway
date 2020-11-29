package parser

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

const (
	TokenFirstLine = iota
	TokenHeader
	TokenBody
)

var _ RequestProvider = &requestFile{}

type requestFile struct {
	path    string
	method  string
	headers map[string]string
	body    strings.Builder
}

func NewRequestFile() *requestFile {
	return &requestFile{
		headers: make(map[string]string),
	}
}

func (r *requestFile) Path(defaultPath string) string {
	if len(r.path) == 0 {
		return defaultPath
	}

	return r.path
}

func (r *requestFile) Method() string {
	return r.method
}

func (r *requestFile) Headers() map[string]string {
	return r.headers
}

func (r *requestFile) Body() io.Reader {
	return strings.NewReader(r.body.String())
}

func (r *requestFile) setHeader(key, val string) {
	r.headers[key] = val
}

func (r *requestFile) addBodyLine(ln string) {
	r.body.WriteString(ln + "\n")
}

func ParseFromReader(reader io.Reader) (RequestProvider, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	res := NewRequestFile()

	token := TokenFirstLine
	for scanner.Scan() {
		ln := scanner.Text()

		if len(ln) > 0 && (ln[0] == '#' || ln[0:1] == "//") {
			continue
		}
		if len(ln) == 0 {
			token = TokenBody
			continue
		}

		switch token {
		case TokenFirstLine:
			var firstLineErr error
			res.method, res.path, firstLineErr = parseFirstLine(ln)
			if firstLineErr != nil {
				return nil, firstLineErr
			}
			token = TokenHeader
			break
		case TokenHeader:
			headerKey, headerVal, headerErr := parseHeaderLine(ln)
			if headerErr != nil {
				return nil, headerErr
			}
			res.setHeader(headerKey, headerVal)
			break
		case TokenBody:
			res.addBodyLine(ln)
			break
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}

	return res, nil
}

func parseFirstLine(line string) (method, path string, err error) {
	split := strings.SplitN(line, " ", 2)

	if len(line) == 0 || len(split) == 0 {
		return "", "", errors.New("empty first header line: " + line)
	}

	if len(split) == 1 {
		method = strings.ToUpper(strings.TrimSpace(split[0]))
		path = ""
	}

	if len(split) == 2 {
		method = strings.ToUpper(strings.TrimSpace(split[0]))
		path = strings.TrimSpace(split[1])
	}

	if len(method) == 0 {
		return "", "", errors.New("key cannot be empty: " + line)
	}

	methodMatched, reErr := regexp.Match("^[a-zA-Z]+$", []byte(method))
	if !methodMatched || reErr != nil {
		return "", "", errors.New("key name contains invalid characters: " + line)
	}

	return
}

func parseHeaderLine(line string) (key, value string, err error) {
	split := strings.SplitN(line, ":", 2)

	if len(split) != 2 {
		return "", "", errors.New("invalid header: " + line)
	}

	key = strings.TrimSpace(split[0])
	if len(key) == 0 {
		return "", "", errors.New("invalid header key: " + line)
	}

	return key, strings.TrimSpace(split[1]), nil
}