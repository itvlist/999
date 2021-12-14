package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"

)

const RETRY_COUNT = 3

func GetURL(id string) string {
	url := "http://news.tvb.com/live/" + id + "?is_hd"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Forwarded-For", "185.225.12.27")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	reg := regexp.MustCompile(`http:[^\"]*\.m3u8`)
	dst := reg.Find(body)
	return string(dst)
}

func InewsHandler(w http.ResponseWriter, r *http.Request) {
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

	var url string
	for i := 0; i < RETRY_COUNT; i++ {
		url = GetURL(id)
		if url != "" {
			break
		}
	}

	if len(url) == 0 {
		http.Error(w, "Cant find m3u8", 503)
		return
	}

	w.Header().Set("Location", url)
	http.Error(w, url, 302)
}
