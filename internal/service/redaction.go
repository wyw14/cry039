package service

import (
	"regexp"
	"strings"
)

var phone = regexp.MustCompile(`1[3-9][0-9]{9}`)

func RedactRemark(v string) string {
	return strings.TrimSpace(phone.ReplaceAllString(v, "[手机号已脱敏]"))
}
