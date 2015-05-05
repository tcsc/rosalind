package gst

import (
	"testing"
	"unicode/utf8"
)

func Test_Init(t *testing.T) {
	u := New("日本語abc日本語abda本語befgda本語beft")
	u.mustBeValid()

	ab := New("abcabxabcd")
	ab.mustBeValid()

	cd := New("cdddcdc")
	cd.mustBeValid()
}

func Test_TreeContainsAllSubstrings(t *testing.T) {
	s := "日本語abc日本語abda本語befgda本語beft"
	tree := New(s)

	for len(s) > 0 {
		_, size := utf8.DecodeRuneInString(s)
		if !tree.Contains(s) {
			t.Errorf("tree should contain %s", s)
		}

		s = s[size:]
	}
}

func decode(s string) []rune {
	result := make([]rune, 0, len(s))
	for _, ch := range s {
		result = append(result, ch)
	}
	return result
}

func Test_SplittingNode(t *testing.T) {
	text := "abcabx"

	parent := newNode(-1)
	n := newNode(0)
	parent.children['a'] = n

	suffixNode := n.split(text, 2, 2)

	if n.children['a'] != suffixNode {
		t.Errorf("Expected new suffix node to be a child of parent")
	}

	if n.str.length != 2 {
		t.Errorf("Expected modified parent to have length (2, got %d", n.str.length)
	}

	if suffixNode.str.length != inf {
		t.Errorf("Expected new suffix to have length (inf), got %d", suffixNode.str.length)
	}
}
