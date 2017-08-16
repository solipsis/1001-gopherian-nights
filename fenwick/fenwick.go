// Package fenwick provides a simple fenwick tree for
// conducting range sum queries on positive values in O(log(n)) time
// TODO: support negative values
package fenwick

type Fenwick []int

// NewFenwick returns a NewFenwick tree which is backed by an
// array with an index for each value in the range of 0 to maxVal
func NewFenwick(maxVal int) Fenwick {
	return make(Fenwick, maxVal+1)
}

// FromList creates a new fenwick tree from a list of starting values.
func FromList(l []int, maxVal int) Fenwick {
	f := make(Fenwick, maxVal+1)
	for _, i := range l {
		f.Adjust(i, 1)
	}
	return f
}

// Adjust increases the cumulative frequency of the given value "v" by the amount "by".
func (f *Fenwick) Adjust(v, by int) {
	for v <= len(*f) {
		(*f)[v] += by
		v += v & -v
	}
}

// QueryRange return the cumulative frequency in range "a" to "b".
func (f *Fenwick) QueryRange(a, b int) int {
	return f.Query(b) - f.Query(a-1)
}

// Query returns the cumulative frequency in range 0 to a.
func (f *Fenwick) Query(a int) int {
	sum := 0
	for a > 0 {
		sum += (*f)[a]
		a -= a & -a
	}
	return sum
}
