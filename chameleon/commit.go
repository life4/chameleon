package chameleon

import (
	"strings"
	"time"
)

type Commits []Commit

func (c Commits) Len() int {
	return len(c)
}

func (c Commits) First() Commit {
	return c[0]
}

func (c Commits) Last() Commit {
	return c[len(c)-1]
}

type Commit struct {
	Hash string
	Time time.Time
	Name string
	Mail string
	Msg  string
	Diff string
}

func ParseCommit(line string) (Commit, error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "|")
	t, err := time.Parse(ISO8601, parts[1])
	if err != nil {
		return Commit{}, err
	}
	c := Commit{
		Hash: parts[0],
		Time: t,
		Name: parts[2],
		Mail: parts[3],
		Msg:  parts[4],
	}
	return c, err
}
