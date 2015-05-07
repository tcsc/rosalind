package gst

import (
	"testing"
	"unicode/utf8"
)

func Test_TreeInvariantsHold(t *testing.T) {
	strings := []string{
		"日本語abc日本語abda本語befgda本語beft",
		"abcabxabcd",
		"cdddcdc",
	}

	for _, s := range strings {
		tree := New(s)
		tree.mustBeValid()
	}
}

func Test_TreeContainsAllSubstrings(t *testing.T) {
	strings := []string{
		"日本語abc日本語abda本語befgda本語beft",
		"abcabxabcd",
		"cdddcdc",
	}

	for _, s := range strings {
		tree := New(s)

		for len(s) > 0 {
			_, size := utf8.DecodeRuneInString(s)
			if !tree.Contains(s) {
				t.Errorf("tree should contain %s", s)
			}

			s = s[size:]
		}
	}
}

func Test_SplittingNode(t *testing.T) {
	text := "abcabx"

	parent := newNode(0, -1)
	n := newNode(0, 0)
	parent.children['a'] = n

	suffixNode := n.split(text, 2)

	if n.str.offset != 0 {
		t.Errorf("Expected parent offset to remain unchanged")
	}

	if n.str.length != 2 {
		t.Errorf("Expected parent to have length (2, got %d", n.str.length)
	}

	if n.children['c'] != suffixNode {
		t.Errorf("Expected new suffix node to be a child of parent")
	}

	if suffixNode.str.offset != 2 {
		t.Errorf("Expected new suffix to have offset 2, got %d",
			suffixNode.str.offset)
	}

	if suffixNode.str.length != inf {
		t.Errorf("Expected new suffix to have length (inf), got %d",
			suffixNode.str.length)
	}
}

func Test_LinkReturnsNextNode(t *testing.T) {
	a := newNode(0, 0)
	b := newNode(0, 42)
	if link(a, b) != b {
		t.Error("Expected link() to return next, but it didn't")
	}
}

func Test_LinkCreatePrevToNext(t *testing.T) {
	a := newNode(0, 0)
	b := newNode(0, 42)
	link(a, b)
	if a.suffix != b {
		t.Error("Expected suffix link to be node b, but it wasn't")
	}
}

func Test_LinkCanTakeNilPrevPtr(t *testing.T) {
	a := newNode(0, 0)
	link(nil, a) // assert this doesn't actually crash
}

func Test_GeneralisedTreeIsValid(t *testing.T) {
	strings := []string{
		"The answer ... is fourty-two!",
		"Fourty-two?",
		"Yes! Fourty-two!",
		"Fourty Two!? We're going to get lynched, aren't we?",
	}

	tree := New(strings...)
	tree.dumpTree("h2g2.dot")
	tree.mustBeValid()
}

func Test_FindAllActuallyFindsAll(t *testing.T) {
	strings := []string{
		"The answer ... is fourty-two!",
		"Fourty-two?",
		"Yes! Fourty-two!",
		"Fourty two!? We're going to get lynched, aren't we?",
	}

	tree := New(strings...)
	points := tree.FindAll("our")

	if len(points) != 4 {
		t.Errorf("Expected %d points, got %d", 4, len(points))
	}

	for _, pt := range points {
		text := tree.Str(pt.Id)[pt.Offset : pt.Offset+3]
		if text != "our" {
			t.Errorf("Expected to find \"our\", got \"%s\"\n", text)
		}
	}
}

func Test_GetStringReturnsOriginalString(t *testing.T) {
	strings := []string{
		"The answer ... is fourty-two!",
		"Fourty-two?",
		"Yes! Fourty-two!",
		"Fourty two!? We're going to get lynched, aren't we?",
	}

	tree := New(strings...)

	for i, s := range strings {
		if tree.Str(i) != s {
			t.Errorf("Expected \"%s\", got \"%s\"", tree.Str(i))
		}
	}
}

func Test_StringIterationReturnsAllStringsInOrder(t *testing.T) {
	strings := []string{
		"The answer ... is fourty-two!",
		"Fourty-two?",
		"Yes! Fourty-two!",
		"Fourty two!? We're going to get lynched, aren't we?",
	}
	tree := New(strings...)

	for s := range tree.Strings() {
		if s != strings[0] {
			t.Errorf("Expected \"%s\", got \"%s\".", strings[0], s)
		}
		strings = strings[1:]
	}
}
