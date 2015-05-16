package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ss := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ss = append(ss, strings.TrimSpace(scanner.Text()))
	}

	if len(ss) < 2 {
		panic("Expected at least 2 strings")
	}

	fmt.Printf("%d\n", HammingDistance(ss[0], ss[1]))
}

func HammingDistance(a, b string) int {
	if len(a) != len(b) {
		panic(fmt.Sprintf("string length mismatch, a: %d, b: %d\n",
			len(a),
			len(b)))
	}

	dh := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			dh++
		}
	}

	return dh
}
