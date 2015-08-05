package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	subs, err := reader.ReadString('\n')
	if err != nil {
		panic(nil)
	}

	text = strings.TrimSpace(text)
	subs = strings.TrimSpace(subs)

	offset := 0
	for {
		i := strings.Index(text, subs)
		if i == -1 {
			break
		}

		hit := offset + i
		offset = hit + 1
		fmt.Printf("%d ", hit+1)
		text = text[i+1:]
	}
}
