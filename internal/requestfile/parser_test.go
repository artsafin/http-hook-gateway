package requestfile

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	TestDefaultPath = "default--path"
)

func TestParseFirstLine(t *testing.T) {
	var tests = []struct {
		input  string
		method string
		path   string
		err    error
	}{
		{
			"POST /some/value?with=query",
			"POST",
			"/some/value?with=query",
			nil,
		},
		{
			"POST    /some/value?with=query",
			"POST",
			"/some/value?with=query",
			nil,
		},
		{
			"GET value",
			"GET",
			"value",
			nil,
		},
		{
			"PUT",
			"PUT",
			"",
			nil,
		},
		{
			"delete",
			"DELETE",
			"",
			nil,
		},
		{
			"",
			"",
			"",
			errors.New("empty first header line: "),
		},
		{
			" ",
			"",
			"",
			errors.New("key cannot be empty:  "),
		},
		{
			" /some/value",
			"",
			"",
			errors.New("key cannot be empty:  /some/value"),
		},
		{
			"/some/value",
			"",
			"",
			errors.New("key name contains invalid characters: /some/value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			method, path, err := parseFirstLine(tt.input)

			if fmt.Sprint(tt.err) != fmt.Sprint(err) {
				t.Errorf("error: got %v, want %v", err, tt.err)
			}
			if method != tt.method {
				t.Errorf("key: got %v, want %v", method, tt.method)
			}
			if path != tt.path {
				t.Errorf("value: got %v, want %v", path, tt.path)
			}
		})
	}
}

func TestParseHeaderLine(t *testing.T) {
	var tests = []struct {
		input string
		key   string
		value string
		err   error
	}{
		{
			"Content-Type: hello: world",
			"Content-Type",
			"hello: world",
			nil,
		},
		{
			"helloworld",
			"",
			"",
			errors.New("invalid header: helloworld"),
		},
		{
			":helloworld",
			"",
			"",
			errors.New("invalid header key: :helloworld"),
		},
		{
			"a:b",
			"a",
			"b",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, value, err := parseHeaderLine(tt.input)

			if fmt.Sprint(tt.err) != fmt.Sprint(err) {
				t.Errorf("error: got %v, want %v", err, tt.err)
			}
			if key != tt.key {
				t.Errorf("key: got %v, want %v", key, tt.key)
			}
			if value != tt.value {
				t.Errorf("value: got %v, want %v", value, tt.value)
			}
		})
	}
}

func TestParseFromReader_Good(t *testing.T) {
	var tests = []struct {
		input   string
		method  string
		path    string
		headers map[string]string
		body    string
	}{
		{
			input: `POST

{"foo": 123}
`,
			method:  "POST",
			path:    TestDefaultPath,
			headers: map[string]string{},
			body:    "{\"foo\": 123}\n",
		},
		{
			input: `POST /some/path

{"foo": 123}
`,
			method:  "POST",
			path:    "/some/path",
			headers: map[string]string{},
			body:    "{\"foo\": 123}\n",
		},
		{
			input: `# Comment on first line
POST /some/path

// This is comment on body
{"foo": 123}
// Last line
`,
			method:  "POST",
			path:    "/some/path",
			headers: map[string]string{},
			body:    "{\"foo\": 123}\n",
		},
		{
			input: `POST
Content-Type: application/json
X-Some-Header: hello

{
	"foo": 123,
	"multiline": "easy"
}`,
			method: "POST",
			path:   TestDefaultPath,
			headers: map[string]string{
				"Content-Type":  "application/json",
				"X-Some-Header": "hello",
			},
			body: `{
	"foo": 123,
	"multiline": "easy"
}
`,
		},
	}

	for index, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			req, err := ParseFromReader(strings.NewReader(tt.input))

			if err != nil {
				t.Errorf("got error for %v: %v", index, err)
				return
			}

			assert.Equal(t, tt.method, req.Method())
			assert.Equal(t, tt.path, req.Path(TestDefaultPath))
			assert.Equal(t, tt.headers, req.Headers())

			bodyBytes, _ := ioutil.ReadAll(req.Body())
			assert.Equal(t, tt.body, string(bodyBytes))
		})
	}
}
