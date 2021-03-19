package chameleon

import "sort"

type Stat struct {
	Path  string
	Count uint32
}

type ViewStat []Stat

func (s ViewStat) Len() int {
	return len(s)
}

func (s ViewStat) Less(i, j int) bool {
	return s[i].Count < s[j].Count
}

func (s ViewStat) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *ViewStat) Sort() {
	sort.Sort(sort.Reverse(s))
}

func (s *ViewStat) Add(path string, count uint32) {
	*s = append(*s, Stat{Path: path, Count: count})
}
