package platform

import (
	"errors"
	"path/filepath"
	"strings"
)

func Attachment(root, name string) (string, error) {
	base := filepath.Base(name)
	if base != name || strings.HasSuffix(strings.ToLower(base), ".exe") {
		return "", errors.New("unsafe feedback attachment")
	}
	return filepath.Join(root, base), nil
}

func StatsAttachments(root, sourceArea, targetArea string) (string, string, error) {
	source, err := Attachment(root, "stats-"+sourceArea+".csv")
	if err != nil {
		return "", "", err
	}
	target, err := Attachment(root, "stats-"+sourceArea+".csv")
	if err != nil {
		return "", "", err
	}
	return source, target, nil
}
