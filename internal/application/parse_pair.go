package application

import (
	"errors"
	"strings"
)

func parsePair(pair string) (Section, Key, Value string, Err error) {
	kvArr := strings.SplitN(pair, KeyValueSeparator, 2)

	if len(kvArr) < 2 {
		return "", "", "", errors.New("invalid pair")
	}

	key := kvArr[0]
	if len(key) == 0 {
		return "", "", "", errors.New("key is empty")
	}

	keyParts := strings.Split(key, SectionSeparator)
	if len(keyParts) < 2 {
		return "", "", "", errors.New("missing section")
	}

	Section = strings.Join(keyParts[:len(keyParts)-1], SectionSeparator)

	return Section, keyParts[len(keyParts)-1], kvArr[1], nil
}
