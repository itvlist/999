package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func SctvHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "http://scgctvshow.sctv.com" + r.URL.Path
	logrus.Debugf("Curl URL:%s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	if strings.HasSuffix(r.URL.Path, ".ts") {
		MyCopy(w, resp.Body)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	content := strings.Replace(string(body), "3-", "1-", -1)
	w.Write([]byte(content))
}
