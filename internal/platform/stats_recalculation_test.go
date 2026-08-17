package platform

import "testing"

func TestStatsRecalculationPipelineAttachments(t *testing.T) {
	source, target, err := StatsAttachments("reports", "a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if source == target {
		t.Fatalf("source and target statistics overwrite the same artifact: %q", source)
	}
}
