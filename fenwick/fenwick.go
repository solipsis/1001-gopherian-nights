package fenwick

type Fenwick []int

func NewFenwick(maxVal int) Fenwick {
	return make(Fenwick, maxVal+1)
}

/*
func FromList(l []int, maxVal int) {
    f := make(Fenwick, maxVal+1)

    make(Fenwick, ma)
}
*/

func (f *Fenwick) adjust(v, by int) {
	for v <= len(*f) {
		(*f)[v] += by
		v += v & -v
	}
}

func (f *Fenwick) QueryRange(a, b int) int {
	return f.Query(b) - f.Query(a-1)
}

func (f *Fenwick) Query(a int) int {
	sum := 0
	for a > 0 {
		sum += (*f)[a]
		a -= a & -a
	}
	return sum
}
