package platform

import (
	"errors"
	"github.com/wyw14/cry039/internal/domain"
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

func UndoReceiptIDs(m domain.Migration, items []domain.Feedback) []string {
	if len(m.AffectedIDs) == 0 {
		return nil
	}
	affected := map[string]struct{}{}
	for _, id := range m.AffectedIDs {
		affected[id] = struct{}{}
	}
	var matched []domain.Feedback
	for _, f := range items {
		if _, ok := affected[f.ID]; ok {
			matched = append(matched, f)
		}
	}
	return domain.IDs(matched)
}
