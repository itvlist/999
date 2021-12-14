package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var host string
var tvfile string

type Result struct {
	Data []Data `json:"data"`
}

type Data struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

var testProgram Data

func PostData(data *Data) {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://"+host+":19394/api/category.json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	time.Sleep(100 * time.Millisecond)
}

func AddURL() {
	fd, err := os.Open(tvfile)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(fd)
	data := Data{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		line_parts := strings.Split(line, ",")
		if len(line_parts) != 1 {
			data.Data += line + "\r\n"
			continue
		}
		if data.Name != "" && data.Data != "" {
			data.Name = base64.StdEncoding.EncodeToString([]byte(data.Name))
			data.Data = base64.StdEncoding.EncodeToString([]byte(data.Data))
			PostData(&data)
			data = Data{}
		}
		data.Name = line
	}
	if data.Name != "" && data.Data != "" {
		data.Name = base64.StdEncoding.EncodeToString([]byte(data.Name))
		data.Data = base64.StdEncoding.EncodeToString([]byte(data.Data))
		PostData(&data)
	}

	if testProgram.Name != "" && testProgram.Data != "" {
		PostData(&testProgram)
	}
	fmt.Println("友窝源添加成功")
}

func DeleteAll() {
	resp, err := http.Get("http://" + host + ":19394/api/categories.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := Result{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	for _, data := range result.Data {
		time.Sleep(100 * time.Millisecond)
		if data.Name == "测试频道" {
			testProgram.Name = base64.StdEncoding.EncodeToString([]byte(data.Name))
			testProgram.Data = base64.StdEncoding.EncodeToString([]byte(data.Data))
		}
		deleteURL := "http://" + host + ":19394/api/delete.json?id=" + data.Id
		res, err := http.Get(deleteURL)
		if err != nil {
			panic(err)
		}
		res.Body.Close()
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <host>\n", os.Args[0])
	os.Exit(0)
}

func main() {
	if len(os.Args) <= 2 {
		Usage()
	}
	host = os.Args[1]
	tvfile = os.Args[2]
	DeleteAll()
	AddURL()
}
