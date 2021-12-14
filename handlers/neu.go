package handlers

import (
	"bytes"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

)

const PROXY_PREFIX = "https://6proxy.6box.cn/index.php?q="

func NeutsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "https://ts.neu6.edu.cn" + strings.TrimPrefix(r.URL.Path, "/neu")
	bt := time.Now()
	resp, err := ProxyGet(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	n, err := io.Copy(w, resp.Body)
	et := time.Now()
	cost := et.Sub(bt).Nanoseconds()
	logrus.Infof("Curl url: %s download: %.2fMB cost: %.2fs speed:%.2fMbps threadNum:%d", url, float64(n)/1024/1024, float64(cost)/1e9, float64(n)/float64(cost)*1e9/1024/1024*8, DefaultThreadNum)
}

func NeuHandler(w http.ResponseWriter, r *http.Request) {
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
	resp, err := ProxyGet(url)
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
	newBody := bytes.Replace(body, []byte("https://ts.neu6.edu.cn"), []byte("http://"+r.Host+"/neu"), -1)
	w.Write(newBody)
}

func ProxyGet(u string) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	u = PROXY_PREFIX + url.QueryEscape(u)
	return client.Get(u)
}
