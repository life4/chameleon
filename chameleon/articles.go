package chameleon

import (
	"sort"
)

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Less(i, j int) bool {
	left := a[i]
	right := a[j]
	c1, _ := left.Commits()
	if c1 == nil {
		return left.Path.Name() < right.Path.Name()
	}
	c2, _ := right.Commits()
	if c1 == nil {
		return left.Path.Name() < right.Path.Name()
	}
	return c1.Created().Time.Unix() < c2.Created().Time.Unix()
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a *Articles) Sort() {
	if a.Len() <= 1 {
		return
	}
	sort.Sort(sort.Reverse(a))
}
