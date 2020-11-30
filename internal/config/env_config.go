package config

import (
	"http-hook-gateway/internal/application"
	"os"
	"strings"
)

const (
	KeyValueSeparator = "="
	SectionSeparator  = "."
)

func LoadFromEnv(prefix string) application.HookMap {
	hookmap := make(application.HookMap)

	normPrefix := strings.Trim(prefix, SectionSeparator) + SectionSeparator

	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, normPrefix) {
			continue
		}
		unprefixedKV := strings.Trim(strings.TrimPrefix(kv, normPrefix), ".")

		section, key, value, err := parsePair(unprefixedKV)
		if err != nil {
			continue
		}

		var ok bool
		var hook *application.HookDef
		if hook, ok = hookmap[section]; !ok {
			hook = &application.HookDef{Name: section}
			hookmap[section] = hook
		}
		assignParamValue(hook, key, value)
	}

	return hookmap
}

func assignParamValue(def *application.HookDef, param, value string) {
	switch param {
	case "accept_url_regex":
		def.AcceptUrlRegex = value
	case "proxy_host":
		def.ProxyHost = value
	case "proxy_path":
		def.ProxyPath = value
	case "request_file":
		def.RequestFile = value
	}
}
