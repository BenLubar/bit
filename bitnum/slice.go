package bitnum

import "sort"

type Numbers []*Number

func (s Numbers) Len() int           { return len(s) }
func (s Numbers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Numbers) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s Numbers) Search(n *Number) int {
	return sort.Search(len(s), func(i int) bool {
		return !s[i].Less(n)
	})
}

func (s Numbers) Contains(n *Number) bool {
	i := s.Search(n)
	return i < len(s) && s[i].Equal(n)
}

func (s *Numbers) Add(n *Number) bool {
	i := s.Search(n)
	if i < len(*s) && (*s)[i].Equal(n) {
		return false
	}

	*s = append(*s, nil)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = n

	return true
}
