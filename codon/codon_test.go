package codon

import (
	"testing"
)

func Test_AllBasesAreCovered(t *testing.T) {
	bases := []byte{'U', 'C', 'A', 'G'}
	for _, a := range bases {
		for _, b := range bases {
			for _, c := range bases {
				_, ok := Table[codon{a, b, c}]
				if !ok {
					t.Errorf("Missing codon for %c%c%c", a, b, c)
				}
			}
		}
	}
}
