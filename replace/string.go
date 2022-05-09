package replace

import "strings"

type stringReplacer struct {
	old string
	new string
}

func NewStringReplacer(old string, new string) Replacer {

	return &stringReplacer{
		old: old,
		new: new,
	}
}

func (r *stringReplacer) Replace(s string) string {
	return strings.ReplaceAll(s, r.old, r.new)
}
