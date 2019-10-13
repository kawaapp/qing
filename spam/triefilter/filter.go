package triefilter

import (
	"github.com/kawaapp/kawaqing/spam"
	"github.com/importcjj/sensitive"
	"strings"
	"sync"
)


// 基于: https://github.com/importcjj/sensitive
type TrieChecker struct {
	sync.Mutex
	*sensitive.Filter
}

func (t *TrieChecker) Validate(text string) (bool, string) {
	return t.Filter.Validate(text)
}

func (t *TrieChecker) Dispose() {
	t.Filter = nil
}

func (t *TrieChecker) setup(str string) {
	words := split(str)
	t.AddWord(words...)
}

func New(words string) spam.SpamChecker {
	trie := &TrieChecker{
		Filter: sensitive.New(),
	}
	if len(words) > 0 {
		trie.setup(words)
	}
	return trie
}

func split(str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool {
		return r == ';' || r == '；'
	})
}