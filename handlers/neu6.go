package handlers

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func Neu6tsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "https://ts.neu6.edu.cn" + r.URL.Path
	bt := time.Now()
	n, err := MultiDownload(w, url, DefaultThreadNum)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	et := time.Now()
	cost := et.Sub(bt).Nanoseconds()
	logrus.Infof("Curl url: %s download: %.2fMB cost: %.2fs speed:%.2fMbps threadNum:%d", url, float64(n)/1024/1024, float64(cost)/1e9, float64(n)/float64(cost)*1e9/1024/1024*8, DefaultThreadNum)
}

func Neu6Handler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	id := r.Form.Get("id")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url := "https://media2.neu6.edu.cn/hls/" + id + ".m3u8"
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
	newBody := bytes.Replace(body, []byte("https://ts.neu6.edu.cn"), []byte("http://"+r.Host), -1)
	w.Write(newBody)
}
