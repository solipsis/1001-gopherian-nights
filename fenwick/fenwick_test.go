package fenwick

import (
	"testing"
)

func TestNew(t *testing.T) {
	f := NewFenwick(1000)
	if len(f) != 1001 {
		t.Errorf("Size should not be '%v'.", len(f))
	}
}

func TestAdjust(t *testing.T) {
	f := NewFenwick(16)
	f.adjust(1, 1)
	// All 2^Nth bits should be set 1/2/4/8/16
	expected := []int{0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1}
	for n := range expected {
		if f[n] != expected[n] {
			t.Errorf("tree should not be %v.", f)
		}
	}
}
