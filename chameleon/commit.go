package chameleon

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05-07:00"

type Commits []Commit

func (c Commits) Len() int {
	return len(c)
}

func (c Commits) Edited() Commit {
	return c[0]
}

func (c Commits) Created() Commit {
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
	if len(parts) < 5 {
		return Commit{}, errors.New("unexpected chunks count")
	}
	t, err := time.Parse(ISO8601, parts[1])
	if err != nil {
		return Commit{}, fmt.Errorf("parse time: %v", err)
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
