package main

import (
	"fmt"
	"github.com/tcsc/rosalind/codon"
	"io/ioutil"
	"os"
)

func main() {
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", translate(string(bytes)))
}

func translate(s string) string {
	n := 0
	c := codon.New()
	result := []rune{}
	for _, x := range s {
		c[n] = byte(x)
		if n == 2 {
			result = append(result, codon.Table[c])
			n = 0
		} else {
			n = n + 1
		}
	}

	return string(result)
}
