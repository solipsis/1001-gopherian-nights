package fenwick

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	f := NewFenwick(1, 1000)
	if len(f.tree) != 1001 {
		t.Errorf("Size should be 1000 not %v", len(f.tree))
	}
}

func TestAdjust(t *testing.T) {
	f := NewFenwick(1, 16)
	f.Adjust(1, 1)
	// All 2^Nth bits (offset by 1) should be set 1/2/4/8/16
	expected := []int{0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1}
	for n := range expected {
		if f.tree[n] != expected[n] {
			t.Errorf("tree should not be %v.", f)
		}
	}
}

func TestFromList(t *testing.T) {
	l := []int{1, 1}
	f := FromList(l, 0, 10)
	fmt.Println(f.tree)
	expected := []int{0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0}
	for n := range expected {
		if f.tree[n] != expected[n] {
			t.Errorf("FromList() expected: %v, actual: %v", expected, f.tree)
			return
		}
	}
}

func TestRangeQuery(t *testing.T) {

	l := []int{2, 4, 5, 5, 6, 6, 6, 7, 7, 8, 9}
	f := FromList(l, 0, 10)

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
