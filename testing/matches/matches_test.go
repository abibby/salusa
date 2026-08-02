package matches

import (
	"testing"
)

func TestEqualTo_Matches(t *testing.T) {
	e := &EqualTo{expected: 5}

	if d := e.Matches(5); d != nil {
		t.Fatalf("expected nil diff for equal values, got %#v", d)
	}

	if d := e.Matches(6); d != nil {
		t.Fatalf("expected nil diff for unequal values, got %#v", d)
	}
}
