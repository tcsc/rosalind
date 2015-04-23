package fasta

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Codon struct {
	Name     string
	Sequence string
}

func ReadFile(filename string) {

}

func Read(reader io.Reader) <-chan Codon {
	ch := make(chan Codon, 2)
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
				// signal error
				return
			}

			line := strings.TrimSpace(s.Text())
			if len(line) == 0 || line[0] == ';' {
				continue
			}

			if line[0] == '>' {
				if seq.Len() > 0 {
					ch <- Codon{name, seq.String()}
				}
				name = line[1:]
				seq.Reset()
			} else {
				seq.WriteString(line)
			}
		}

		if seq.Len() > 0 {
			ch <- Codon{name, seq.String()}
		}
	}()

	return ch
}
