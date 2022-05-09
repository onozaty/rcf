package replace

import (
	"regexp"

	"github.com/pkg/errors"
)

type regexpReplacer struct {
	regex       *regexp.Regexp
	replacement string
}

func NewRegexpReplacer(regexStr string, replacement string) (Replacer, error) {

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &regexpReplacer{
		regex:       regex,
		replacement: replacement,
	}, nil
}

func (r *regexpReplacer) Replace(s string) string {
	return r.regex.ReplaceAllString(s, r.replacement)
}
