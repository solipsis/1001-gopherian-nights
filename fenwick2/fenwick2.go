// Package fenwick2 provides an improved fenwick tree for
// conducting range sum queries in O(log(n)) time
// this implementation has been upgraded to support negative values
// and ranges other than 1 to n
package fenwick2

// Fenwick tree represented by an int array and an offset
// that is used to support negative numbers and ranges outside of 0 to n
type Fenwick struct {
	tree   []int
	offset int
}

// NewFenwick returns a NewFenwick tree which is backed by an
// array with an index for each value in the range of min to max inclusive
func NewFenwick(min, max int) Fenwick {
	f := Fenwick{
		tree: make([]int, max-min+2),
	}
	var offset = 0
	if min < 1 {
		offset = (-min) + 1
	} else if min > 1 {
		offset = -min + 1
	}
	f.offset = offset

	return f
}

// FromList creates a new fenwick tree from a list of starting values.
func FromList(l []int, min, max int) Fenwick {
	f := NewFenwick(min, max)
	for _, i := range l {
		f.Adjust(i, 1)
	}
	return f
}

// Adjust increases the cumulative frequency of the given value "v" by the amount "by".
func (f *Fenwick) Adjust(v, by int) {
	v += f.offset
	for v < len((*f).tree) {
		(*f).tree[v] += by
		v += v & -v
	}
}

// QueryRange return the cumulative frequency in range "a" to "b".
func (f *Fenwick) QueryRange(a, b int) int {
	return f.Query(b) - f.Query(a-1)
}

// Query returns the cumulative frequency in range min to a.
func (f *Fenwick) Query(a int) int {
	a += f.offset
	sum := 0
	for a > 0 {
		sum += (*f).tree[a]
		a -= a & -a
	}
	return sum
}
