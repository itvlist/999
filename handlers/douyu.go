package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type DouyuResult struct {
	Errno int       `json:"error"`
	Data  DouyuData `json:"data"`
}

type DouyuData struct {
	URL string `json:"hls_url"`
}

func DouyuHandler(w http.ResponseWriter, r *http.Request) {
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
	url := "https://m.douyu.com/html5/live?roomId=" + id
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
	result := DouyuResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	dst := result.Data.URL
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
