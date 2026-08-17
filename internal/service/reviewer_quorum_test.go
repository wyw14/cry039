package service

import "testing"

func TestReviewerQuorumPipelineService(t *testing.T) {
	if DistinctReviewers([]string{"Alice", " alice "}) {
		t.Fatal("equivalent reviewer identities were treated as distinct")
	}
	if !DistinctReviewers([]string{"Alice", "Bob"}) {
		t.Fatal("two different reviewers should be accepted")
	}
}
