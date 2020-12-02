package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_loadFromEnvPairs(t *testing.T) {
	type args struct {
		prefix string
		pairs  []string
	}
	tests := []struct {
		name string
		args args
		want HookMap
	}{
		{
			name: "simple",
			args: args{
				prefix: "test",
				pairs: []string{
					"test.section1.accept_url_regex=value1",
					"test.section1.proxy_host=value2",
					"test.section1.proxy_path=value3",
					"test.section1.request_file=value4",
					"test.section2.accept_url_regex=value5",
					"test.section2.proxy_host=value6",
					"test.section2.proxy_path=value7",
					"test.section2.request_file=value8",

					"test.section2.accept_url_regex=value5 overriden",
				},
			},
			want: map[string]*HookDef{
				"section1": {
					Name:           "section1",
					AcceptUrlRegex: "value1",
					ProxyHost:      "value2",
					ProxyPath:      "value3",
					RequestFile:    "value4",
				},
				"section2": {
					Name:           "section2",
					AcceptUrlRegex: "value5 overriden",
					ProxyHost:      "value6",
					ProxyPath:      "value7",
					RequestFile:    "value8",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadFromEnvPairs(tt.args.prefix, tt.args.pairs)
			assert.Equal(t, tt.want, got)
		})
	}
}
