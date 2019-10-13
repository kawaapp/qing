package triefilter

import "testing"

func TestTrieChecker_Validate(t *testing.T) {
	filter := New("习包子;去死")
	ok, _ := filter.Validate("一去二三里")
	if !ok {
		t.Error("it should works, no spam")
	}
	ok, w := filter.Validate("学习包子")
	if ok || w != "习包子" {
		t.Error("warning! not find spam words")
	}
}

func TestNoWordsProvided(t *testing.T)  {
	filter := New("")
	ok, _ := filter.Validate("abc")
	if !ok {
		t.Error("it should works, if no spam words provided")
	}
}

func TestSplit(t *testing.T) {
	array := split("A;B；C")
	if len(array) != 3 {
		t.Error("string split size err, get:", len(array))
	}
	if array[0] != "A" || array[1] != "B" || array[2] != "C" {
		t.Error("string splite err, get:", array[0], array[1], array[2])
	}
}