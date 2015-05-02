package gst

import (
	"testing"
)

func Test_Init(t *testing.T) {
	New("abcabxabcd")
	New("cdddcdc")
}

func decode(s string) []rune {
	result := make([]rune, 0, len(s))
	for _, ch := range s {
		result = append(result, ch)
	}
	return result
}

func Test_SplittingNode(t *testing.T) {
	text := decode("abcabx")

	parent := newNode(-1, nil)
	n := newNode(0, parent)
	parent.children['a'] = n

	newParent, nn, newChild := n.split(text, 2, 2)

	if parent.children['a'] != newParent {
		t.Errorf("Expected new parent to be a child or parent")
	}

	if nn != n {
		t.Errorf("Expected nn to be the original node")
	}

	_ = newChild

	if newParent.length != 2 {
		t.Errorf("Expected new parent to have length (2, got %d", newParent.length)
	}

	if newChild.length != inf {
		t.Errorf("Expected new child to have length (inf), got %d", newChild.length)
	}
}
