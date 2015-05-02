package fasta

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func Test_ReadNonExistantFileProducessError(t *testing.T) {
	var count = 0
	for codon := range ReadFile("/no/such/file") {
		count++
		if codon.Error == nil {
			t.Fatal("Expected an error")
		}
	}

	if count != 1 {
		t.Fatalf("Expeced to hit loop body only once, gog %d", count)
	}
}

func Test_ReadingAnEmptyStringReturnsNoRecords(t *testing.T) {
	count := 0
	buf := bytes.NewBufferString("")
	for _ = range Read(buf) {
		count++
	}

	if count != 0 {
		t.Errorf("Expected no records, got %d.", count)
	}
}

func Test_ReadingASingleRecordWorks(t *testing.T) {

	buf := bytes.NewBufferString(">SomeName\nGAT\nTACA")
	count := 0
	for codon := range Read(buf) {
		if codon.Name != "SomeName" {
			t.Errorf("Expected name \"SomeName\", got \"%s\"", codon.Name)
		}

		if codon.Sequence != "GATTACA" {
			t.Errorf("Expected sequence \"GATTACA\", got \"%s\"", codon.Name)
		}
	}

	if count != 0 {
		t.Errorf("Expected count 1, got %d.", count)
	}
}

func Test_TrailingNewlinesAreIgnored(t *testing.T) {

	buf := bytes.NewBufferString(">SomeName\nGAT\nTACA\r\n")
	count := 0
	for codon := range Read(buf) {
		if codon.Name != "SomeName" {
			t.Errorf("Expected name \"SomeName\", got \"%s\"", codon.Name)
		}

		if codon.Sequence != "GATTACA" {
			t.Errorf("Expected sequence \"GATTACA\", got \"%s\"", codon.Name)
		}
	}

	if count != 0 {
		t.Errorf("Expected count 1, got %d.", count)
	}
}

func Test_CommentLinesAreIgnored(t *testing.T) {
	buf := bytes.NewBufferString(">SomeName\nGAT\nTACA\n;>AnotherName\nGATAAAATTTTAACACACCA")
	count := 0
	for codon := range Read(buf) {
		var expected Codon
		switch count {
		case 0:
			expected = Codon{"SomeName", "GATTACAGATAAAATTTTAACACACCA", nil}
			break

		default:
			t.Fatal("Too many records returned")
		}

		if codon != expected {
			t.Fatalf("Expected %#v, got %#v", expected, codon)
		}
		count++
	}
}

func Test_MultipleRecordsAreReadInOrder(t *testing.T) {
	buf := bytes.NewBufferString(">SomeName\nGAT\nTACA\n>AnotherName\nGATAAAATTTTAACACACCA")
	count := 0
	for codon := range Read(buf) {
		var expected Codon
		switch count {
		case 0:
			expected = Codon{"SomeName", "GATTACA", nil}
			break

		case 1:
			expected = Codon{"AnotherName", "GATAAAATTTTAACACACCA", nil}
			break

		default:
			t.Fatal("Too many records returned")
		}

		if codon != expected {
			t.Fatalf("Expected %#v, got %#v", expected, codon)
		}
		count++
	}
}

type signallingReader struct {
	reader io.Reader
	signal chan bool
}

func (self *signallingReader) Read(buf []byte) (int, error) {
	return self.reader.Read(buf)
}

func (self *signallingReader) Close() error {
	self.signal <- true
	return nil
}

func Test_ReadCallsCloseOnClosableThings(t *testing.T) {
	signal := make(chan bool)
	reader := signallingReader{
		reader: bytes.NewBufferString(">SomeName\nGAT\nTACA\r\n"),
		signal: signal}

	for _ = range Read(&reader) {
	}

	close_called := false
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case <-signal:
		close_called = true

	case <-timeout:
		t.Fatal("Timed out waiting for closed signal")
	}

	if !close_called {
		t.Error("Close was not called")
	}
}
