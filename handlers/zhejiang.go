//
package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

var reRegx = regexp.MustCompilePOSIX(`([^#]+\.ts)`)

func ZhejiangApiHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url := ""

	id := r.Form.Get("id")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url = fmt.Sprintf("http://hw-m-l.cztv.com/channels/lantian/%s/1080p.m3u8", id)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	url = fmt.Sprintf("http://hw-m-l.cztv.com/%s", r.URL.Path[4:])
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	newbody := reRegx.ReplaceAll(body, []byte("/zj/channels/lantian/"+id+"/$0"))
	//	w.Header().Set("Host", "hw-m-l.cztv.com")
	//	w.Header().Set("Origin", "http://hw-m-l.cztv.com")
	//	w.Header().Set("Referer",  fmt.Sprintf("http://hw-m-l.cztv.com/channels/lantian/%s/", id))
	logrus.Info(string(newbody))
	w.Write(newbody)
}
func ZhejiangHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	url := fmt.Sprintf("http://hw-m-l.cztv.com/%s", r.URL.Path[4:])

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()

	MyCopy(w, resp.Body)

}
