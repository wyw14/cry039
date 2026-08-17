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
