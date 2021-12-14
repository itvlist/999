package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

)

func YoutubeApiHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
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

func YoutubeVideoHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
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
	logrus.Debugf("TS URL:%s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-type", "application/octet-stream")
	MyCopy(w, resp.Body)
}

func YoutubeIndexHandler(w http.ResponseWriter, r *http.Request) {
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
	q := r.Form.Get("q")
	url := "https://www.youtube.com/watch?v=" + id
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
	re := regexp.MustCompile(`hlsvp.*m3u8`)
	hls := re.Find(body)
	if len(hls) < 8 {
		http.Error(w, "Cant't find m3u8 url", 503)
		return
	}
	dst := string(hls[8:])
	dst = strings.Replace(dst, "\\/", "/", -1)
	if q == "" {
		dst = strings.Replace(dst, "https://manifest.googlevideo.com", "http://"+r.Host, 1)
		w.Header().Set("Location", dst)
		http.Error(w, dst, 302)
		return
	}
	seq, err := strconv.Atoi(q)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	resp2, err := http.Get(dst)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp2.Body.Close()
	body, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	re = regexp.MustCompile(`https://.*\.m3u8`)
	urls := re.FindAllString(string(body), -1)
	seq = len(urls) - seq
	if seq < 0 || seq >= len(urls) {
		errStr := fmt.Sprintf("Quality not exist, total: %d", len(urls))
		http.Error(w, errStr, 503)
		return
	}

	dst = strings.Replace(urls[seq], "https://manifest.googlevideo.com", "http://"+r.Host, 1)
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
