package main

import (
	"testing"
)

func Test_lcs1(t *testing.T) {
	a := "ABCDEFXNARFO"
	b := "XBCDYYFNBARFX"
	rval := findLCS(a, b)

	if len(rval) != 2 {
		t.Fail()
	}

	if rval[0] != "BCD" {
		t.Errorf("%s :-(", rval)
	}
}

func Test_lcs2(t *testing.T) {
	a := "ABCDEFXNARFO"
	b := "CDEF"
	rval := findLCS(a, b)

	if len(rval) != 1 {
		t.Fail()
	}

	if rval[0] != "CDEF" {
		t.Errorf("%s :-(", rval)
	}
}
