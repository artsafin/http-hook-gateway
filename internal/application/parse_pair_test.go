package application

import (
	"errors"
	"fmt"
	"testing"
)

func TestParsePair(t *testing.T) {
	var tests = []struct {
		input   string
		section string
		key     string
		value   string
		err     error
	}{
		{
			"section.key=value",
			"section",
			"key",
			"value",
			nil,
		},
		{
			"section.with.many.dots.key=value",
			"section.with.many.dots",
			"key",
			"value",
			nil,
		},
		{
			"some_section.some_key=value=with=equals",
			"some_section",
			"some_key",
			"value=with=equals",
			nil,
		},
		{
			"some_section.some_key==some_value",
			"some_section",
			"some_key",
			"=some_value",
			nil,
		},
		{
			"some_section..some_key=some_value",
			"some_section.",
			"some_key",
			"some_value",
			nil,
		},
		{
			"some_key=some_value",
			"",
			"",
			"",
			errors.New("missing section"),
		},
		{
			"=some_value",
			"",
			"",
			"",
			errors.New("key is empty"),
		},
		{
			"",
			"",
			"",
			"",
			errors.New("invalid pair"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			section, key, value, err := parsePair(tt.input)

			if fmt.Sprint(tt.err) != fmt.Sprint(err) {
				t.Errorf("error: got %v, want %v", err, tt.err)
			}
			if section != tt.section {
				t.Errorf("section: got %v, want %v", section, tt.section)
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
