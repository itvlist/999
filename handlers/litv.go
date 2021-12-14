package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"

)

func LitvHandler(w http.ResponseWriter, r *http.Request) {
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
	url := "http://btsu4k5-hisng.cdn.hinet.net/live/pool/" + id + "/litv-pc/index.m3u8"
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
	re := regexp.MustCompile(`.*.m3u8`)
	hls := re.FindAll(body, -1)
	if len(hls) < 1 {
		http.Error(w, "Cant't find m3u8 url", 503)
		return
	}
	dst := "http://btsu4k5-hisng.cdn.hinet.net/live/pool/" + id + "/litv-pc/" + string(hls[len(hls)-1])
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
