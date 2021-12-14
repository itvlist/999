package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func QmHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := "http://liveal.quanmin.tv" + r.URL.Path
	logrus.Debugf("Curl url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	MyCopy(w, resp.Body)
}
