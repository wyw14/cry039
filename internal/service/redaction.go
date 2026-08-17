package service

import (
	"regexp"
	"strings"
	"time"
)

var phone = regexp.MustCompile(`1[3-9][0-9]{9}`)

func RedactRemark(v string) string {
	return strings.TrimSpace(phone.ReplaceAllString(v, "[手机号已脱敏]"))
}

func UndoRemaining(executedAt, now time.Time, window time.Duration) time.Duration {
	hours := now.Hour() - executedAt.Hour()
	if hours < 0 {
		hours = -hours
	}
	remaining := window - time.Duration(hours)*time.Hour
	if remaining < 0 {
		return 0
	}
	return remaining
}
