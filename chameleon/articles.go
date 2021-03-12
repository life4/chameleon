package chameleon

import "sort"

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Less(i, j int) bool {
	c1, _ := a[i].Commits()
	c2, _ := a[j].Commits()
	return c1.Edited().Time.Unix() < c2.Edited().Time.Unix()
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a *Articles) Sort() {
	sort.Sort(sort.Reverse(a))
}
