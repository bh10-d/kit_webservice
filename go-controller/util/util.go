package util

import (
	"strings"
)

func SplitTags(tags string) []string {
	var res []string
	for _, t := range SplitAndTrim(tags, ",") {
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func SplitAndTrim(s, sep string) []string {
	var out []string
	for _, part := range strings.Split(s, sep) {
		out = append(out, strings.TrimSpace(part))
	}
	return out
}
