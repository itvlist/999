package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		line_parts := strings.Split(line, ",")
		url := line_parts[len(line_parts)-1]
		if !strings.HasPrefix(url, "http") {
			fmt.Printf("%s\n", line)
			continue
		}
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("%s\ttimeout\n", line)
			continue
		}
		fmt.Printf("%s\t%d\n", line, resp.StatusCode)
		resp.Body.Close()
		time.Sleep(1 * time.Second)
	}
}
