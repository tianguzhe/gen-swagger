package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func LowerCamelCase(s string) string {
	return _camelCase(s, true)
}

func UpperCamelCase(s string) string {
	return _camelCase(s, false)
}

func _camelCase(s string, hasLower bool) string {
	var strBuild strings.Builder

	tc := cases.Title(language.English)

	for k, v := range strings.Split(s, "-") {
		if k == 0 {
			if hasLower {
				strBuild.WriteString(v)
			} else {
				strBuild.WriteString(tc.String(v))
			}
		} else {
			strBuild.WriteString(tc.String(v))
		}
	}

	return strBuild.String()
}
