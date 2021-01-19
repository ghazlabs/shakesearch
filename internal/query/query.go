package query

import "strings"

type Query string

func (q Query) GetWords() []string {
	// maybe we need to improve this
	return strings.Split(string(q), " ")
}
