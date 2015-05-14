package main

import (
	"fmt"
	"github.com/tcsc/rosalind/fasta"
	"os"
)

func gcContent(s string) float64 {
	count := 0
	for _, ch := range s {
		if ch == 'G' || ch == 'C' {
			count++
		}
	}
	return (float64(count) / float64(len(s))) * 100.0
}

func main() {
	gcMax := 0.0
	leader := ""
	for str := range fasta.ReadFile(os.Args[1]) {
		if str.Error != nil {
			panic(str.Error)
		}

		gc := gcContent(str.Sequence)
		if gc > gcMax {
			gcMax = gc
			leader = str.Name
		}
	}

	fmt.Printf("%s\n%f\n", leader, gcMax)
}
