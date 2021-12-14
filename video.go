package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
)

func videoHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	host := r.Form.Get("hls_chunk_host")
	if host == "" {
		re := regexp.MustCompile(`[\w\d-]+\.googlevideo.com`)
		host = re.FindString(r.URL.Path)
		if host == "" {
			http.Error(w, http.StatusText(503), 503)
			return
		}
	}
	url := "https://" + host + r.URL.EscapedPath()
	if len(r.URL.RawQuery) > 0 {
		url += "?" + r.URL.RawQuery
	}
	logrus.Debug("TS URL:%s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-type", "application/octet-stream")
	io.Copy(w, resp.Body)
}
