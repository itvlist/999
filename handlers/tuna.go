package handlers

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

func TunaTsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "https://iptv.tsinghua.edu.cn" + strings.TrimPrefix(r.URL.Path, "/tuna")
	bt := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	n, err := io.Copy(w, resp.Body)
	et := time.Now()
	cost := et.Sub(bt).Nanoseconds()
	logrus.Infof("Curl url: %s download: %.2fMB cost: %.2fs speed:%.2fMbps", url, float64(n)/1024/1024, float64(cost)/1e9, float64(n)/float64(cost)*1e9/1024/1024*8)
}

func TunaHandler(w http.ResponseWriter, r *http.Request) {
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
	url := "https://iptv.tsinghua.edu.cn/hls/" + id + ".m3u8"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasSuffix(line, ".ts") {
			line = "http://" + r.Host + "/tuna/hls/" + line
		}
		fmt.Fprintln(w, line)
	}
}
