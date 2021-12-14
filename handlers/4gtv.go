package handlers

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func FourgtvTsHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Request URL", r.URL)
	url := "https://4gtvfreemobile-cds.cdn.hinet.net" + r.URL.Path
	logrus.Debug("Curl url", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	MyCopy(w, resp.Body)
}

func FourgtvApiHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "https://p-sirona-yond-4gtv.svc.litv.tv" + r.URL.Path
	logrus.Debug("Curl url", url)
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
	body = bytes.Replace(body, []byte("https://4gtvfreemobile-cds.cdn.hinet.net"), []byte("http://"+r.Host), -1)
	w.Write(body)
}

func FourgtvHandler(w http.ResponseWriter, r *http.Request) {
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
	jsonStr := `{"jsonrpc":"2.0","id":12,"method":"LoadService.GetURLsNoAuth","params":{"AssetId":"` + id + `","DeviceType":"mobile","MediaType":"channel"}}`
	req, err := http.NewRequest("POST", "https://twproxy02.svc.litv.tv/cdi/4gtv/rpc", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := gclient.Do(req)
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
	re := regexp.MustCompile(`http[^"]*.m3u8`)
	hls := re.Find(body)
	if len(hls) < 1 {
		http.Error(w, "Cant't find m3u8 url", 503)
		return
	}
	resp2, err := http.Get(string(hls))
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp2.Body.Close()
	body, err = ioutil.ReadAll(resp2.Body)
	re = regexp.MustCompile(`.*\.m3u8`)
	m3u8 := re.FindAll(body, -1)
	if len(m3u8) < 1 {
		http.Error(w, "Cant't find m3u8 url", 503)
		return
	}
	dst := string(bytes.Replace(hls, []byte("https://p-sirona-yond-4gtv.svc.litv.tv"), []byte("http://"+r.Host), -1))
	dst = strings.Replace(dst, "master.m3u8", string(m3u8[len(m3u8)-1]), -1)
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
