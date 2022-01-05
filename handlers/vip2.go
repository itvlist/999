package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var baseAPI = "http://42.193.18.62:9999/analysis.php?v="

func Vip2Handler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	URL := r.Form.Get("url")
	if URL == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	resp, err := http.Get(baseAPI + url.PathEscape(URL))
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`[\s\S]+var\surls\s=\s"(.+)";[\s\S]+`)
	rawHTML := string(bytes)
	if matches := re.FindStringSubmatch(rawHTML); len(matches) >= 2 {
		URL = matches[1]
	}
	parsedURL, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}
	extIndex := strings.LastIndexByte(parsedURL.Path, '.')
	if extIndex != -1 {
		URLExt := strings.ToLower(parsedURL.Path[extIndex+1:])
		switch URLExt {
		case "mp4", "mkv", "avi":
			// No need to download m3u8
			//	return URL, URLExt, false
		}
	}
	//return URL, "", true

	http.RedirectHandler(URL, 302).ServeHTTP(w, r)

}
