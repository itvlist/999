package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		line_parts := strings.Split(line, ",")
		if len(line_parts) != 2 {
			continue
		}
		fmt.Printf("#EXTINF:0,%s\n", line_parts[0])
		fmt.Printf("%s\n", line_parts[1])
	}
}
