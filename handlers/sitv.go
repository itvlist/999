package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func getIP() (string, error) {
	resp, err := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}`)
	ip := re.Find(body)
	return string(ip), nil
}

func SitvHandler(w http.ResponseWriter, r *http.Request) {
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
	ct := r.Form.Get("ct")
	if ct == "" {
		ct = "1"
	}
	ip, err := getIP()
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	url := "http://www.sitv.com.cn/GetPlayPath/GetPlayPath?type=LIVE&se=sitv&ct=" + ct + "&code=" + id + "&ip=" + ip
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
	dst := strings.Replace(string(body), "2300000", "4000000", 1)
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
