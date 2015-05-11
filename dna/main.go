package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	totals := [4]int{0, 0, 0, 0}
	stdin := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := stdin.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		switch ch {
		case 'A':
			totals[0]++
		case 'C':
			totals[1]++
		case 'G':
			totals[2]++
		case 'T':
			totals[3]++
		}
	}
	fmt.Printf("%d %d %d %d", totals[0], totals[1], totals[2], totals[3])
}
