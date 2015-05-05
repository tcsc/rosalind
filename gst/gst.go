package gst

import (
	"fmt"
	"os"
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
/// the text itself, it only stores pointers to the text in an string.
type node struct {
	suffix   *node
	str      substring
	children map[rune]*node
}

/// newNode creates and initialises a node in its defauts state: starting at a
/// given offset and extending for the remainder or the internal text
func newNode(start int) *node {
	return &node{
		str:      substring{index: 0, offset: start, length: inf},
		suffix:   nil,
		children: make(map[rune]*node),
	}
}

/// split splits a node into two parts
func (self *node) split(text string, i, length int) *node {
	childLength := inf
	if self.str.length != inf {
		childLength = self.str.length - length
	}

	newChild := &node{
		str: substring{
			index:  0,
			offset: self.str.offset + length,
			length: childLength},
		suffix:   nil,
		children: self.children,
	}

	self.str.length = length
	key := decodeRune(text, newChild.str.offset)
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

/// label() generates a label for the node, representing its text value
func (self *node) label(text string) string {
	if self.str.offset == -1 {
		return "root"
	}

	last := len(text)
	if self.str.length != -1 {
		last = self.str.offset + self.str.length
	}
	if self.str.offset >= len(text) {
		panic("oob")
	}
	return string(text[self.str.offset:last])
}

func (self *node) isShorterThan(n int) bool {
	if self.str.length != inf {
		return self.str.offset+self.str.length <= n
	}
	return false
}

///
type SuffixTree struct {
	root   *node
	corpus []string
}

/// Creates a new suffix treen and initialises it from the supplied string.
func New(s string) SuffixTree {
	tree := SuffixTree{
		root:   newNode(-1),
		corpus: []string{s},
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
func (self *activePointState) slide(child *node, index int, text string) bool {
	if child.str.length != inf && self.length >= child.str.length {
		self.length -= child.str.length
		self.edge = decodeRune(text, index-self.length)
		self.node = child
		fmt.Printf("Sliding to %s[%c]\n", child.label(text), self.edge)
		return true
	}
	return false
}

/// Generates a suffix link between a the nodes iff prevNode is not nil.
func link(prev, next *node, s string) *node {
	if prev != nil {
		fmt.Printf("Making suffix link from %s to %s\n",
			prev.label(s),
			next.label(s))
		prev.suffix = next
	}
	return next
}

func decodeRune(text string, offset int) rune {
	r, _ := utf8.DecodeRuneInString(text[offset:])
	return r
}

func (self *node) format(text string) string {
	result := fmt.Sprintf("%s {", self.label(text))
	for k, v := range self.children {
		result = result + fmt.Sprintf("%c: %s, ", k, v.label(text))
	}
	result += "}"
	return result
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
/// Based on code from http://pastie.org/5925812#72-106
func (self *SuffixTree) Insert(s string) { //, index int) {
	active := activePointState{self.root, '\x00', 0}
	remainder := 0

	fmt.Printf("len(s): %d\n", len(s))

	i := 0
	text := s
	for len(text) > 0 {
		c, charlen := utf8.DecodeRuneInString(text)
		remainder++
		prefix := s[:i+charlen]

		var prevNode *node = nil

		fmt.Println("========")
		fmt.Printf("#%d: %c (%d bytes)\n", i, c, charlen)

		for remainder > 0 {
			if active.length == 0 {
				active.edge = c
			}

			fmt.Printf("\nActive Node: %s\nActive Length: %d\nRemainder: %d\nActive edge %c\n",
				active.node.format(prefix),
				active.length,
				remainder,
				active.edge)

			activeChild, ok := active.node.children[active.edge]
			if !ok {
				fmt.Printf("New node for %c\n", c)
				newChild := newNode(i)
				active.node.children[active.edge] = newChild

				fmt.Printf("Active Node is now %s\n", active.node.format(prefix))
				prevNode = newChild //link(prevNode, newChild, prefix)
			} else {
				fmt.Printf("Active edge length: %d\n", activeChild.str.length)
				if active.slide(activeChild, i, s) {
					fmt.Printf("Sliding\n")
					continue
				}

				fmt.Printf("Decoding s[%d + %d = %d]\n",
					activeChild.str.offset,
					active.length,
					activeChild.str.offset+active.length)
				fmt.Printf("Looking for %c, found %c\n",
					c,
					decodeRune(s, activeChild.str.offset+active.length))

				if decodeRune(s, activeChild.str.offset+active.length) == c {
					active.length += charlen
					//prevNode = link(prevNode, active.node, prefix)
					break
				} else {
					fmt.Printf("Splitting %s\n", activeChild.format(prefix))
					grandChild := activeChild.split(s, i, active.length)

					newChild := newNode(i)
					activeChild.children[c] = newChild

					fmt.Printf(
						"Parent: %s\nSplit suffix: %s\nNew child: %s\n",
						activeChild.format(prefix),
						grandChild.format(prefix),
						newChild.format(prefix))

					prevNode = link(prevNode, newChild, prefix)
				}
			}
			remainder--

			if active.node == self.root && active.length > 0 {
				x, n := utf8.DecodeRuneInString(s[i-active.length:])
				fmt.Printf("Old leader: %c\n", x)

				active.length -= n
				active.edge = decodeRune(s, (i - active.length))
				fmt.Printf("new leader: %c\n", active.edge)
			} else {
				if active.node.suffix != nil {
					fmt.Printf("Following suffix link to %s\n",
						active.node.suffix.label(prefix))
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

/// Contains checks to see if the tree contains a given substring
func (self *SuffixTree) Contains(s string) bool {
	node := self.root
	nodeStr := ""
	for _, ch := range s {
		if len(nodeStr) == 0 {
			if n, ok := node.children[ch]; !ok {
				return false
			} else {
				node = n
				nodeStr = self.nodeString(node)
			}
		}
		otherChar, size := utf8.DecodeRuneInString(nodeStr)
		if ch != otherChar {
			return false
		}

		nodeStr = nodeStr[size:]
	}
	return true
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

	queue := []*node{self.root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		label := n.label(self.corpus[n.str.index])
		file.WriteString(fmt.Sprintf("\"%p\" [label=\"%s\"]\n", n, label))
		for k, v := range n.children {
			if k == '\x00' {
				k = '?'
			}
			file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [label=\"%c\"]\n", n, v, k))
		}

		if n.suffix != nil {
			file.WriteString(fmt.Sprintf("\"%p\" -> \"%p\" [style=\"dotted\"]\n", n, n.suffix))
		}

		queue = append(queue, n.childNodes()...)
	}
}
