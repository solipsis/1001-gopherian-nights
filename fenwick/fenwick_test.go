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
	f.Adjust(1, 1)
	// All 2^Nth bits should be set 1/2/4/8/16
	expected := []int{0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1}
	for n := range expected {
		if f[n] != expected[n] {
			t.Errorf("tree should not be %v.", f)
		}
	}
}

func TestFromList(t *testing.T) {
	l := []int{2, 4, 5, 5, 6, 6, 6, 7, 7, 8, 9}
	f := FromList(l, 10)
	expected := []int{0, 0, 1, 0, 2, 2, 5, 2, 10, 1, 1}
	for n := range expected {
		if f[n] != expected[n] {
			t.Errorf("FromList() expected: %v, actual: %v", expected, f)
		}
	}
}

func TestRangeQuery(t *testing.T) {

	l := []int{2, 4, 5, 5, 6, 6, 6, 7, 7, 8, 9}
	f := FromList(l, 10)

	type rangeQueryTests struct {
		a, b     int
		expected int
	}
	var rangeTests = []rangeQueryTests{
		{1, 1, 0},
		{1, 2, 1},
		{1, 6, 7},
		{1, 10, 11},
		{3, 6, 6},
	}
	for _, rt := range rangeTests {
		actual := f.QueryRange(rt.a, rt.b)
		if actual != rt.expected {
			t.Errorf("RangeQuery(%d,%d): expected %d, actual %d", rt.a, rt.b, rt.expected, actual)
		}
	}
}
