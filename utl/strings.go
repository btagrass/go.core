package utl

import (
	"strings"
)

// 替换
func Replace(s string, oldnews ...string) string {
	replacer := strings.NewReplacer(oldnews...)

	return replacer.Replace(s)
}

// 分割
func Split(s string, seps ...rune) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, sep := range seps {
			if sep == r {
				return true
			}
		}

		return false
	})
}
