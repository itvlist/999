package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Request URL:%s", r.URL)
	url := "https://manifest.googlevideo.com" + r.URL.EscapedPath()
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	content := strings.Replace(string(body), "https://manifest.googlevideo.com", "http://"+r.Host, -1)
	re := regexp.MustCompile(`https://[\w\d-]+\.googlevideo.com`)
	content = re.ReplaceAllString(content, "http://"+r.Host)
	w.Header().Set("Content-type", "application/vnd.apple.mpegurl")
	w.Write([]byte(content))
}
