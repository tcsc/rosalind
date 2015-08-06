package main

import (
	"fmt"
	"github.com/tcsc/rosalind/fasta"
	"os"
)

var bases = []uint8{'A', 'C', 'G', 'T'}

type profile map[uint8][]int

func main() {
	matrix := []string{}

	// load the matrix
	for s := range fasta.ReadFile(os.Args[1]) {
		if s.Error != nil {
			panic(s.Error)
		}
		matrix = append(matrix, s.Sequence)
	}

	if len(matrix) == 0 {
		return
	}

	p := buildProfile(matrix)
	c := buildConsensus(p)

	for _, b := range c {
		fmt.Printf("%c", b)
	}
	fmt.Printf("\n")

	printProfile(p)
}

func newProfile(n int) profile {
	result := make(profile)
	for _, c := range bases {
		result[c] = make([]int, n)
	}
	return result
}

func buildProfile(matrix []string) profile {
	n := len(matrix[0])
	p := newProfile(n)
	for i := 0; i < n; i++ {
		for _, base := range bases {
			count := 0
			for _, seq := range matrix {
				if seq[i] == base {
					count++
				}
			}
			p[base][i] = count
		}
	}
	return p
}

func buildConsensus(p profile) []uint8 {
	n := len(p['A'])
	consensus := make([]uint8, n)
	for i := 0; i < n; i++ {
		max := -1
		for _, base := range bases {
			count := p[base][i]
			if count > max {
				consensus[i] = base
				max = count
			}
		}
	}
	return consensus
}

func printProfile(p profile) {
	for _, base := range bases {
		fmt.Printf("%c:", base)
		for _, n := range p[base] {
			fmt.Printf(" %d", n)
		}
		fmt.Printf("\n")
	}
}
