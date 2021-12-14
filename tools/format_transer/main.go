package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Name string
	Urls []string
}

func NewItem() *Item {
	return &Item{Urls: make([]string, 1)}
}

func (item *Item) Add(name, url string) {
	if name == item.Name {
		item.Urls = append(item.Urls, url)
		return
	}
	item.Print()
	item.Reset(name, url)
}

func (item *Item) Print() {
	if item.Name != "" && len(item.Urls) != 0 {
		fmt.Printf("%s,%s\n", item.Name, strings.Join(item.Urls, "#"))
	}
	item.Clear()
}

func (item *Item) Clear() {
	item.Name = ""
	item.Urls = make([]string, 1)
}

func (item *Item) Reset(name, url string) {
	item.Name = name
	item.Urls = []string{url}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-r" {
		FormatR()
		return
	}
	Format()
}

func Format() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line_parts := strings.Split(sc.Text(), ",")
		if len(line_parts) != 2 {
			fmt.Println(sc.Text())
			continue
		}
		urls := strings.Split(line_parts[1], "#")
		for _, url := range urls {
			fmt.Printf("%s,%s\n", line_parts[0], url)
		}
	}
}

func FormatR() {
	item := NewItem()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line_parts := strings.Split(sc.Text(), ",")
		if len(line_parts) != 2 {
			item.Print()
			fmt.Println(sc.Text())
			continue
		}
		item.Add(line_parts[0], line_parts[1])
	}
	item.Print()
}
