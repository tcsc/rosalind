package gst

import (
	"fmt"
	"os"
)

const (
	/// The length value indicating "as long as the string"
	inf = -1
)

/// node represents a chunk of text inside the suffix tree. It doesn't store
/// the text itself, it only stores pointers to the text in an string.
type node struct {
	parent   *node
	suffix   *node
	offset   int
	length   int
	children map[rune]*node
}

/// newNode creates and initialises a node in its defauts state: starting at a
/// given offset and extending for the remainder or the internal text
func newNode(start int, parent *node) *node {
	return &node{
		parent:   parent,
		offset:   start,
		suffix:   nil,
		length:   inf,
		children: make(map[rune]*node),
	}
}

/// split splits a node into two parts
func (self *node) split(text []rune, i, length int) (*node, *node, *node) {
	oldParent := self.parent
	newLeaf := newNode(i, nil)
	newParent := &node{
		parent: oldParent,
		suffix: self.suffix,
		offset: self.offset,
		length: length,
		children: map[rune]*node{
			text[self.offset+1]: self,
			text[i]:             newLeaf},
	}

	newLeaf.parent = newParent

	oldParent.children[text[self.offset]] = newParent
	self.parent = newParent
	self.offset += length
	if self.length != inf {
		self.length -= length
	}
	self.suffix = nil

	return newParent, self, newLeaf
}

/// A key/value pair use when iterating over the children of a node
type child struct {
	index rune
	n     *node
}

/// childNodes returns a slice of key/value pairs, representing the
/// child suffixes of the node, and the indices used to address them.
func (self *node) childNodes() []child {
	result := make([]child, len(self.children))
	i := 0
	for k, v := range self.children {
		result[i] = child{index: k, n: v}
		i++
	}
	return result
}

/// id generates an ID string for the node.
func (self *node) id() string {
	return fmt.Sprintf("%p", self)
}

/// label() generates a label for the node, representing its text value
func (self *node) label(text []rune) string {
	last := len(text)
	if self.length != -1 {
		last = self.offset + self.length
	}
	if self.offset >= len(text) {
		panic("oob")
	}
	return string(text[self.offset:last])
}

func (self *node) isShorterThan(n int) bool {
	if self.length != inf {
		return self.offset+self.length <= n
	}
	return false
}

///
type SuffixTree struct {
	root *node
	text []rune
}

/// Creates a new suffix treen and initialises it from the supplied string.
func New(s string) SuffixTree {
	tree := SuffixTree{
		root: newNode(-1, nil),
		text: make([]rune, 0),
	}
	tree.Insert(s)
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
func (self *activePointState) slide(child *node, text []rune) bool {
	if child.length != inf && self.length >= child.length {
		self.length -= child.length
		self.edge = text[child.offset]
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

/// Insert inserts a new string into the suffix tree.
func (self *SuffixTree) Insert(s string) {
	active := activePointState{self.root, '\x00', 0}
	remainder := 0
	for _, c := range s {
		i := len(self.text)
		self.text = append(self.text, c)
		remainder++

		var prevNode *node = nil

		for remainder > 0 {
			if active.length == 0 {
				active.edge = c
			}

			activeChild, ok := active.node.children[active.edge]
			if !ok {
				newChild := newNode(i, active.node)
				active.node.children[active.edge] = newChild
				prevNode = link(prevNode, newChild)
			} else {
				if active.slide(activeChild, self.text) {
					continue
				}

				if self.text[activeChild.offset+active.length] == c {
					active.length++
					prevNode = link(prevNode, active.node)
					break
				} else {
					np, _, _ := activeChild.split(self.text, i, active.length)
					prevNode = link(prevNode, np)
				}
			}
			remainder--

			if active.node == self.root && active.length > 0 {
				active.length--
				active.edge = self.text[i-active.length]
			} else {
				if active.node.suffix != nil {
					active.node = active.node.suffix
				} else {
					active.node = self.root
				}
			}
		}
	}
}

/// dumpTree writes the tree out to a dot-formatted file  for diagnstic
/// purposes.
func (self *SuffixTree) dumpTree(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("Failed to open graph output file")
	}
	defer file.Close()

	file.WriteString("digraph G {\n")
	defer file.WriteString("}")

	file.WriteString(fmt.Sprintf("\"%p\" [label=\"root\"]\n", self.root))

	queue := self.root.childNodes()
	for len(queue) > 0 {
		n := queue[0].n
		x := queue[0].index
		if x == '\x00' {
			x = '?'
		}
		queue = queue[1:]

		label := n.label(self.text)
		file.WriteString(fmt.Sprintf("\"%p\" [label=\"%s\"]\n", n, label))
		file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [label=\"%c\"]\n", n.parent, n, byte(x)))

		if n.suffix != nil {
			file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [style=\"dotted\"]\n", n, n.suffix))
		}

		queue = append(queue, n.childNodes()...)
	}
}
