package application

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"http-hook-gateway/internal/requestfile"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var testJsonData, _ = getJsonBodyData(strings.NewReader(`{"id": 1, "type": "tall", "details": {"color": "black", "used": true}, "perks": ["clever", "tall"]}`))

var testJsonInputReq = RequestSummary{
	Method:     "GET",
	RemoteAddr: "100.101.102.103",
	Headers: http.Header{
		"Content-Type":    []string{"application/json"},
		"X-Custom-Header": []string{"header value 1", "header value 2"},
	},
	Body:         testJsonData,
	Scheme:       "https",
	User:         "user",
	UserPassword: "pass",
	Hostname:     "example.org",
	Path:         "/initial/path",
	Fragment:     "fragment",
	Query: url.Values{
		"query1": []string{"value1"},
		"query2": []string{"value2"},
	},
}

func TestInterpolateRequestfile(t *testing.T) {
	type args struct {
		src  requestfile.RequestFile
		data *RequestSummary
	}
	tests := []struct {
		args args
		want requestfile.RequestFile
		err  error
	}{
		{
			args: args{
				src: requestfile.NewDto(
					"{{ .Method }}",
					`/proxied/path/initial/path/?{{ range $key, $value := .Query -}}
						param_{{ $key }}=param_{{ index $value 0 }}&
					{{- end }}{{ query . "query1" }}#fragment-fragment`,
					map[string]string{
						"X-Custom-Header": `Some Value and {{ headervalues . "X-Custom-Header" }}`,
						"Content-Type":    `{{ header . "Content-Type" }}`,
						"X-New-Header":    `some new value`,
					},
					`{"foo": "{{ .Body.type }}", "bar": {{ .Body.details.used }}, "baz": {{ json .Body.perks }} }`,
				),
				data: &testJsonInputReq,
			},
			want: requestfile.NewDto(
				"GET",
				`/proxied/path/initial/path/?param_query1=param_value1&param_query2=param_value2&value1#fragment-fragment`,
				map[string]string{
					"X-Custom-Header": `Some Value and header value 1; header value 2`,
					"Content-Type":    `application/json`,
					"X-New-Header":    `some new value`,
				},
				`{"foo": "tall", "bar": true, "baz": ["clever","tall"] }`,
			),
			err: nil,
		},
	}
	for ti, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", ti), func(t *testing.T) {
			got, err := interpolateRequestfile(tt.args.src, tt.args.data)

			if fmt.Sprint(tt.err) != fmt.Sprint(err) {
				t.Errorf("error: got `%v`, want `%v`", err, tt.err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
