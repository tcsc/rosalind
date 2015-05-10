package main

import (
	"fmt"
	"github.com/tcsc/rosalind/fasta"
	"os"
	"time"
)

type matrix struct {
	arr []int
	m   int
	n   int
}

func newMatrix(m int, n int) matrix {
	return matrix{
		arr: make([]int, m*n),
		m:   m,
		n:   n,
	}
}

func (self *matrix) get(i, j int) int {
	return self.arr[(i*self.n)+j]
}

func (self *matrix) set(i, j, val int) {
	self.arr[(i*self.n)+j] = val
}

func findLCS(a, b string) []string {

	matrix := newMatrix(len(a), len(b))
	max := 0
	rval := make([]string, 0)

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				length := 0
				if i == 0 || j == 0 {
					length = 1
				} else {
					length = matrix.get(i-1, j-1) + 1
				}

				text := a[i-length+1 : i+1]

				if len(text) > 0 {
					if length > max {
						max = length
						rval = []string{text}
					} else if length == max {
						rval = append(rval, text)
					}
				}

				matrix.set(i, j, length)
			}
		}
	}

	return rval
}

func main() {
	fmt.Printf("Loading %s...\n", os.Args[1])

	init := true
	lcs := map[string]bool{}
	start := time.Now()
	for c := range fasta.ReadFile(os.Args[1]) {
		if c.Error != nil {
			panic(c.Error)
		}

		if init {
			lcs = map[string]bool{c.Sequence: true}
			init = false
		} else {
			tmp := map[string]bool{}
			for s, _ := range lcs {
				for _, newStr := range findLCS(c.Sequence, s) {
					tmp[newStr] = true
				}
			}
			lcs = tmp
		}
	}

	fmt.Printf("Computation took %s\n", time.Since(start))

	for s, _ := range lcs {
		fmt.Println(s)
	}
}
