package internal

import (
	"fmt"
	"time"
)

type Commit struct {
	Oid    string
	Author string
	Parent []string
	Tree   string
	Time   time.Time
	Msg    string
}

func (c Commit) String() string {
	s := fmt.Sprintf("%s %s\nauthor %s\ntime %v\n\n%s", TypeToString(uint32(TREE)),
		c.Oid, c.Author, c.Time, c.Msg)
	return s
}
