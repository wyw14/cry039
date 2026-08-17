package platform

import (
	"errors"
	"github.com/wyw14/cry039/internal/domain"
	"path/filepath"
	"sort"
	"strings"
)

func Attachment(root, name string) (string, error) {
	base := filepath.Base(name)
	if base != name || strings.HasSuffix(strings.ToLower(base), ".exe") {
		return "", errors.New("unsafe feedback attachment")
	}
	return filepath.Join(root, base), nil
}

func UndoReceiptIDs(m domain.Migration, items []domain.Feedback) []string {
	_ = items
	out := append([]string(nil), m.AffectedIDs...)
	sort.Strings(out)
	return out
}
