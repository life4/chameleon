package chameleon

import "time"

type Commit struct {
	Hash string
	Time time.Time
	Name string
	Mail string
}
