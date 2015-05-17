package main

import (
	"fmt"
)

func main() {
	k := 0.0 // Homozyguous Dominant
	m := 0.0 // Heterozygous
	n := 0.0 // Homozygous Recessive

	_, err := fmt.Scanf("%f %f %f", &k, &m, &n)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%f\n", prob(k, m, n))
}

func prob(k, m, n float64) float64 {
	/*
	 *  Probability of dominance:
	 *
	 *  K    K     (K/n) * (K-1/n-1)
	 *  K    M     (K/n) * (M/n-1)
	 *  K    N     (K/n) * (N/n-1)
	 *
	 *  M    K     (M/n) * (K/n-1)
	 *  M    M     (M/n) * (M-1/n-1) * 0.75
	 *  M    N     (M/n) * (N/n-1) * 0.5
	 *
	 *  N    K     (N/n) * (K/n-1)
	 *  N    M     (N/n) * (M/n-1) * 0.5
	 *  N    N     0
	 */

	// summing all these together we get
	t := k + m + n

	pk := (k * (k - 1)) + (k * m) + (k * n)
	pm := (m * k) + (m * (m - 1) * 0.75) + (m * n * 0.5)
	pn := (n * k) + (n * m * 0.5)
	den := t * (t - 1)

	return (pk + pm + pn) / den
}
