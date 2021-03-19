package chameleon

import "sort"

type Stat struct {
	Path  string
	Count uint32
	Repo  Repository
}

func (s Stat) URLs() URLs {
	return URLs{
		Repository: s.Repo,
		Path:       Path(s.Path),
	}
}

type ViewStat struct {
	Stats []Stat
	Max   uint32
}

func (s ViewStat) Len() int {
	return len(s.Stats)
}

func (s ViewStat) Less(i, j int) bool {
	return s.Stats[i].Count < s.Stats[j].Count
}

func (s ViewStat) Swap(i, j int) {
	s.Stats[i], s.Stats[j] = s.Stats[j], s.Stats[i]
}

func (s *ViewStat) Sort() {
	sort.Sort(sort.Reverse(s))
}

func (s *ViewStat) Add(path string, count uint32) {
	s.Stats = append(s.Stats, Stat{Path: path, Count: count})
	if count > s.Max {
		s.Max = count
	}
}

func (stats ViewStat) SetRepo(repo Repository) {
	for i, stat := range stats.Stats {
		stat.Repo = repo
		stats.Stats[i] = stat
	}
}
