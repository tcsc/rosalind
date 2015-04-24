package main

import (
	"fmt"
	"os"
	"strconv"
)

//   | A  | A  ||   |  A |  A ||   | A  | A  ||   | A  | a  ||   | A  | a  ||
// A | AA | AA || A | AA | AA || a | Aa | Aa || A | AA | Aa || a | Aa | aa ||
// A | AA | AA || a | Aa | Aa || a | Aa | Aa || a | Aa | aa || a | Aa | aa ||

func main() {
	probs := [...]float64{
		2.0, // AA-AA
		2.0, // AA-Aa
		2.0, // AA-aa
		1.5, // Aa-Aa
		1.0, // Aa-aa
		0.0}

	var k [6]float64
	for i, s := range os.Args[1:] {
		f, _ := strconv.ParseFloat(s, 64)
		k[i] = f
	}

	sum := 0.0
	for i, p := range probs {
		sum += k[i] * p
	}

	fmt.Printf("Expected offspring with dominant phenotype: %f", sum)
}
