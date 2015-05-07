package gst

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	/// The length value indicating "as long as the string"
	inf = -1
)

/// Represents a chunk of text stored as a UTF-8 string.
type substring struct {
	index  int
	offset int
	length int
}

/// node represents a chunk of text inside the suffix tree. It doesn't store
/// the text itself, it only stores pointers to the text in an external string.
type node struct {
	suffix   *node
	str      substring
	children map[rune]*node
}

/// newNode creates and initialises a node in its defauts state: starting at a
/// given offset and extending for the remainder or the internal text
func newNode(stringId, start int) *node {
	return &node{
		str:      substring{index: stringId, offset: start, length: inf},
		suffix:   nil,
		children: make(map[rune]*node),
	}
}

/// split splits a node into two parts
func (self *node) split(text string, length int) *node {
	childLength := inf
	if self.str.length != inf {
		childLength = self.str.length - length
	}

	newChild := &node{
		str: substring{
			index:  self.str.index,
			offset: self.str.offset + length,
			length: childLength},
		suffix:   nil,
		children: self.children,
	}

	self.str.length = length
	key := decodeRune(text, length)
	self.children = map[rune]*node{key: newChild}

	return newChild
}

/// childNodes returns a slice of key/value pairs, representing the
/// child suffixes of the node, and the indices used to address them.
func (self *node) childNodes() []*node {
	result := make([]*node, len(self.children))
	i := 0
	for _, v := range self.children {
		result[i] = v
		i++
	}
	return result
}

/// id generates an ID string for the node.
func (self *node) id() string {
	return fmt.Sprintf("%p", self)
}

func (self *node) isLeaf() bool {
	return len(self.children) == 0
}

///
type SuffixTree struct {
	root   *node
	corpus []string
}

/// Creates a new suffix treen and initialises it from the supplied string.
func New(strings ...string) SuffixTree {
	tree := SuffixTree{
		root:   newNode(-1, -1),
		corpus: make([]string, 0, 1),
	}

	for _, s := range strings {
		tree.Insert(s)
	}
	return tree
}

/// activePointState defines a struct for managing the current insertion point
type activePointState struct {
	node   *node
	edge   rune
	length int
}

/// edgeTarget fetches a pointer to the currently active child node, i.e. the
/// child of the currently active node pointed to by the active edge. Returns
/// nil if no edge is active, or no such child exists
func (self *activePointState) edgeTarget() *node {
	if result, ok := self.node.children[self.edge]; ok {
		return result
	}

	if self.edge != '\x00' {
		panic("We're missing a child node!")
	}

	return nil
}

/// slide moves the active point along a link to the next child node, if it is
/// appropriate to do so. Returns true if the active point has benn modified,
/// false if it has been left unchanged.
func (self *activePointState) slide(child *node, index int, text string) bool {
	if child.str.length != inf && self.length >= child.str.length {
		self.length -= child.str.length
		self.edge = decodeRune(text, index-self.length)
		self.node = child
		return true
	}
	return false
}

/// Generates a suffix link between a the nodes iff prevNode is not nil.
func link(prev, next *node) *node {
	if prev != nil {
		prev.suffix = next
	}
	return next
}

/// Decodes a rune starting at a given offset inside a string
func decodeRune(text string, offset int) rune {
	r, _ := utf8.DecodeRuneInString(text[offset:])
	return r
}

/// Asserts the invariants of a completed tree
func (self *SuffixTree) mustBeValid() {
	queue := []*node{self.root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for k, v := range n.children {
			// assert that each node's key is the leading character of the
			// string
			if k != self.nodeChar(v, 0) {
				panic(fmt.Sprintf("Expected index %c for node %s. Got %c",
					self.nodeChar(v, 0),
					self.nodeString(v),
					k))
			}
			queue = append(queue, v.childNodes()...)
		}
	}
}

/// nodeChar fetches the i'th character in the substring represented by the
/// node. Asking for a character outside the substring range will result in
/// undefined behaviour.
func (self *SuffixTree) nodeChar(n *node, i int) rune {
	s := self.corpus[n.str.index]
	ch, _ := utf8.DecodeRuneInString(s[n.str.offset+i:])
	return ch
}

/// nodeChar fetches the substring represented by the node. Asking for a
/// character outside the substring range will result in undefined behaviour.
func (self *SuffixTree) nodeString(n *node) string {
	s := self.corpus[n.str.index]
	if n.str.length == inf {
		return s[n.str.offset:]
	} else {
		return s[n.str.offset : n.str.offset+n.str.length]
	}
}

/// Insert inserts a new string into the suffix tree.
func (self *SuffixTree) Insert(s string) {
	id := len(self.corpus)
	taggedText := fmt.Sprintf("%s\x00%08x", s, id)
	self.corpus = append(self.corpus, taggedText)
	self.index(id)
}

func (self *SuffixTree) split(n *node, i int) *node {
	return n.split(self.nodeString(n), i)
}

/// Indexes a string in the corpus
/// Based on code from http://pastie.org/5925812#72-106
func (self *SuffixTree) index(index int) { //, index int) {
	active := activePointState{self.root, '\x00', 0}
	remainder := 0

	i := 0
	text := self.corpus[index]
	str := text
	for len(text) > 0 {
		c, charlen := utf8.DecodeRuneInString(text)
		remainder++
		var prevNode *node = nil

		for remainder > 0 {
			// if we're not already tracking a branch of the active node...
			if active.length == 0 {
				active.edge = c
			}

			// look up the active branch.
			activeChild, ok := active.node.children[active.edge]
			if !ok {
				// branch does not exist - better start it!
				newChild := newNode(index, i)
				active.node.children[active.edge] = newChild
				prevNode = link(prevNode, active.node)
			} else {
				// if we have reached the end of the active branc, it's time to
				// move down the tree to the branch's target node
				if active.slide(activeChild, i, str) {
					// ... and try the current suffix again
					continue
				}

				// look at the character at the insertion point, does it match?
				if self.nodeChar(activeChild, active.length) == c {
					// yep - we can just keep tracking this branch as it already
					// contains the current suffix.
					active.length += charlen
					prevNode = link(prevNode, active.node)
					break
				} else {
					// nope - we need to split the active node at the insertion
					// point so we can insert a new node that encodes our active
					// suffix

					// fmt.Printf("\nSplitting: \"%s\"\n", self.nodeString(activeChild))

					self.split(activeChild, active.length)
					newChild := newNode(index, i)
					activeChild.children[c] = newChild

					// fmt.Printf("prefix:    \"%s\"\n", self.nodeString(activeChild))
					// fmt.Printf("suffix:    \"%s\"\n", self.nodeString(gch))
					// fmt.Printf("new child: \"%s\"\n", self.nodeString(newChild))

					prevNode = link(prevNode, activeChild)
				}
			}
			remainder--

			if active.node == self.root && active.length > 0 {
				_, n := utf8.DecodeRuneInString(str[i-active.length:])
				active.length -= n
				active.edge = decodeRune(str, (i - active.length))
			} else {
				if active.node.suffix != nil {
					active.node = active.node.suffix
				} else {
					active.node = self.root
				}
			}
		}

		i += charlen
		text = text[charlen:]
	}
}

/// find() Searches through the tree to find a given pattern. If the pattern
/// exists, find returns the node and offset that indicates the *end* of the
/// pattern. Returns (nil, -1) if the pattern can't be found.
func (self *SuffixTree) find(s string) (*node, int) {
	node := self.root
	nodeStr := ""
	index := 0
	for _, ch := range s {
		if len(nodeStr) == 0 {
			if n, ok := node.children[ch]; !ok {
				return nil, 0
			} else {
				node = n
				index = 0
				nodeStr = self.nodeString(node)
			}
		}
		otherChar, size := utf8.DecodeRuneInString(nodeStr)
		if ch != otherChar {
			return nil, 0
		}
		index += size
		nodeStr = nodeStr[size:]
	}
	return node, index
}

/// Contains checks to see if the tree contains a given substring
func (self *SuffixTree) Contains(s string) bool {
	n, _ := self.find(s)
	return n != nil
}

type StringLoc struct {
	Id     int
	Offset int
}

func (self *SuffixTree) Str(i int) string {
	return self.corpus[i]
}

func (self *SuffixTree) FindAll(s string) []StringLoc {
	type point struct {
		n      *node
		length int
	}

	result := make([]StringLoc, 0)
	n, offset := self.find(s)
	if n == nil {
		return result
	}

	q := []point{point{n: n, length: n.str.length - offset}}
	var pt point
	for len(q) > 0 {
		pt, q = q[len(q)-1], q[:len(q)-1]
		if pt.n.isLeaf() {
			str := self.corpus[pt.n.str.index]
			offset := len(str) - pt.length - len(s)
			loc := StringLoc{
				Id:     pt.n.str.index,
				Offset: offset,
			}
			result = append(result, loc)
		} else {
			for _, child := range pt.n.children {
				nextPoint := point{
					n:      child,
					length: pt.length + len(self.nodeString(child)),
				}
				q = append(q, nextPoint)
			}
		}
	}

	return result
}

/// dumpTree writes the tree out to a dot-formatted file for diagnostic
/// purposes.
func (self *SuffixTree) dumpTree(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("Failed to open graph output file")
	}
	defer file.Close()

	file.WriteString("digraph G {\n")
	defer file.WriteString("}")

	queue := []*node{self.root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		label := ""
		if n.str.index < 0 {
			label = "root"
		} else {
			label = strings.Replace(self.nodeString(n), "\x00", "(null)", -1)
		}

		file.WriteString(fmt.Sprintf("\"%p\" [label=\"'%s'\"]\n", n, label))
		for k, v := range n.children {
			if k == '\x00' {
				k = '?'
			}
			file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [label=\"'%c'\"]\n", n, v, k))
		}

		if n.suffix != nil {
			file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [style=\"dotted\"]\n", n, n.suffix))
		}

		queue = append(queue, n.childNodes()...)
	}
}
