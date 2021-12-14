package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var gCurlTimeMap = make(map[string]int64)

func ByrApiHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "http://tv.byr.cn:8888" + r.URL.Path
	logrus.Debug("Curl url", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	if strings.HasSuffix(r.URL.Path, "index.m3u8") {
		go CurlCount(r.URL.Path)
	}
	MyCopy(w, resp.Body)
}

func CurlCount(path string) {
	ts := gCurlTimeMap[path]
	now := time.Now().Unix()
	if now-ts >= 25 {
		gCurlTimeMap[path] = now
		url := fmt.Sprintf("http://tv.byr.cn/player_count/res.gif?play_url=http://tv.byr.cn:8888%s&refer=http://tv.byr.cn/tv-show&title=BYR-IPTV", path)
		logrus.Debugf("Curl count url:%s", url)
		resp, err := http.Get(url)
		if err != nil {
			logrus.Warnf("curl count url: %s error: %s", url, err)
		} else {
			resp.Body.Close()
		}
		if len(gCurlTimeMap) > 1000 {
			gCurlTimeMap = make(map[string]int64)
		}
	}
}

func ByrHandler(w http.ResponseWriter, r *http.Request) {
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
	url := "http://tv.byr.cn/tv-show-detail/" + id
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
	re := regexp.MustCompile(`http://.*/index.m3u8`)
	hls := re.Find(body)
	if len(hls) < 8 {
		http.Error(w, "Cant't find m3u8 url", 503)
		return
	}
	dst := strings.Replace(string(hls), "http://tv.byr.cn:8888", "http://"+r.Host, 1)
	w.Header().Set("Location", dst)
	http.Error(w, http.StatusText(302), 302)
}
