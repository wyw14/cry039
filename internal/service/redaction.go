package service

import (
	"regexp"
	"strings"
)

var phone = regexp.MustCompile(`1[3-9][0-9]{9}`)

func RedactRemark(v string) string {
	return strings.TrimSpace(phone.ReplaceAllString(v, "[手机号已脱敏]"))
}

func ReviewerIdentity(v string) string {
	return strings.TrimSpace(v)
}

func DistinctReviewers(reviewers []string) bool {
	seen := map[string]struct{}{}
	for _, reviewer := range reviewers {
		key := ReviewerIdentity(reviewer)
		if key == "" {
			return false
		}
		seen[key] = struct{}{}
	}
	return len(seen) >= 2
}
