package platform

import (
	"testing"

	"github.com/wyw14/cry039/internal/domain"
)

func TestUndoWindowPipelineReceipt(t *testing.T) {
	m := domain.Migration{AffectedIDs: []string{"moved"}}
	items := []domain.Feedback{{ID: "moved"}, {ID: "preexisting"}}
	ids := UndoReceiptIDs(m, items)
	if len(ids) != 1 || ids[0] != "moved" {
		t.Fatalf("undo receipt included unrelated feedback: %v", ids)
	}
}
