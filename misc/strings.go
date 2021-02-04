package misc

import (
	uuid "github.com/google/uuid"
	"sort"
	"strings"
)

type StringSet map[string]bool

func NewUuid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// string unifier
type StringUnifier struct {
	dest       []string
	dupChecker StringSet
	sorted     bool
}

func NewStringUnifier() *StringUnifier {
	uni := new(StringUnifier)
	uni.dest = make([]string, 0)
	uni.dupChecker = make(StringSet)
	return uni
}

func (self *StringUnifier) Add(item string) bool {
	if _, ok := self.dupChecker[item]; !ok {
		self.dest = append(self.dest, item)
		self.dupChecker[item] = true
		self.sorted = false
		return true
	}
	return false
}

func (self *StringUnifier) Result() []string {
	if !self.sorted {
		// sort dest
		sort.Slice(self.dest, func(i, j int) bool {
			return self.dest[i] < self.dest[j]
		})
		self.sorted = true
	}
	return self.dest

}

// InList
func StringInList(a string, candidates ...string) bool {
	for _, ca := range candidates {
		if ca == a {
			return true
		}
	}
	return false
}
