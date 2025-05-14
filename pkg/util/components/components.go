package components

import (
	"regexp"

	"github.com/moonkit02/dearer/pkg/util/normalize_key"
	"github.com/moonkit02/dearer/pkg/util/regex"
)

var keyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bbucket\b`),
	regexp.MustCompile(`\bstore\b`),
	regexp.MustCompile(`\bstorage\b`),
}

func KeyIsRelevant(key string) bool {
	return regex.AnyMatch(keyPatterns, normalize_key.Normalize(key))
}
