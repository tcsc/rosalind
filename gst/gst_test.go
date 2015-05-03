package gst

import (
	"testing"
)

func Test_Init(t *testing.T) {
	u := New("日本語abc日本語abda本語befg")
	u.dumpTree("unicode.dot")
	//New("abcabxabcd")
	//New("cdddcdc")
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

	parent := newNode(-1, nil)
	n := newNode(0, parent)
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
