package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := stdin.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if ch == 'T' {
			ch = 'U'
		}
		fmt.Printf("%c", ch)
	}
}
