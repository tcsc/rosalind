package basestring

import (
	"bytes"
	"testing"
)

func Test_CreatingBaseStringWithValidCharsSucceeds(t *testing.T) {
	s, err := FromString("GATTACA")
	if err != nil {
		t.Fatalf("Conversion failed: %s", err.Error())
	}
	e := []byte{0x21, 0x33, 0x42, 0x02}

	if s.Length() != 7 {
		t.Errorf("Expected length == 7, got %d", s.Length())
	}
	if bytes.Compare(s.chars, e) != 0 {
		t.Errorf("expected %#v, got %#v", e, s.chars)
	}
}

func Test_CreatingBaseStringWithInvalidCharsFails(t *testing.T) {
	_, err := FromString("NarfZortTroz")
	if err == nil {
		t.Fatal("Expected conversion to fail")
	}
}
