package main

// lineStatistic prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]

	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "lineStatistic: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}

	for line, n := range counts {
		fmt.Printf("%d\t%s\n", n, line)
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
