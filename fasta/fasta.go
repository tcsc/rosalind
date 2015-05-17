package fasta

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type String struct {
	Name     string
	Sequence string
	Error    error
}

func ReadFile(filename string) <-chan String {
	file, err := os.Open(filename)
	if err != nil {
		ch := make(chan String, 1)
		ch <- String{"", "", err}
		close(ch)
		return ch
	}
	return Read(file)
}

func Read(reader io.Reader) <-chan String {
	ch := make(chan String, 2)
	go func() {
		// When we exit, close the input stream if its closeable
		defer func() {
			if closer, ok := reader.(io.Closer); ok {
				closer.Close()
			}
		}()

		// When we exit, close the channel back to the caller
		defer close(ch)

		name := ""
		var seq bytes.Buffer
		s := bufio.NewScanner(reader)
		for s.Scan() {
			if s.Err() != nil {
				ch <- String{"", "", s.Err()}
				return
			}

			line := strings.TrimSpace(s.Text())
			if len(line) == 0 || line[0] == ';' {
				continue
			}

			if line[0] == '>' {
				if seq.Len() > 0 {
					ch <- String{name, seq.String(), nil}
				}
				name = line[1:]
				seq.Reset()
			} else {
				seq.WriteString(line)
			}
		}

		if seq.Len() > 0 {
			ch <- String{name, seq.String(), nil}
		}
	}()

	return ch
}
