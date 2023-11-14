package utils

import "testing"

func TestRandSeq(t *testing.T) {
	n := 100
	set := make(map[string]bool)
	for i := 0; i < n; i++ {
		s := RandSeq(6)
		if len(s) != 6 {
			t.Errorf("RandSeq(%d) returned string of length %d, want %d", 6, len(s), 6)
		}
		set[s] = true
	}

	if len(set) != n {
		t.Errorf("RandSeq(%d) returned %d unique strings, want %d", 6, len(set), n)
	}
}
