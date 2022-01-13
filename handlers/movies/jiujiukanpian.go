package movies

import (
	"fmt"
	"github.com/iawia002/annie/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
//https://www.jiujiukanpian.com/
func init()  {
	MoiveStationMap["https://www.jiujiukanpian.com"] = &MovieStation{
		Name: "久久影视网",
		HostUrl: "https://www.jiujiukanpian.com",
		ContentType: "application/vnd.apple.mpegurl",
		Referer: "https://www.jiujiukanpian.com/",
		RootPath: true,
		ReRegxp: regexp.MustCompilePOSIX("([^#]+\\.ts)"),
		UrlBuildFunc: func(requestUrl string, stationInfo MovieStation) string {
			req, err := http.NewRequest("GET", requestUrl, nil)
			if err != nil {
				return ""
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			req.Header.Set("accept", `*/*`)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}
			id := utils.MatchOneOf(string(body), `/play/\?url=[^"]+`)[0]
			if id == "" {
				return ""
			}
			t := id[strings.Index(id, "?") + 1:]
			params := strings.Split(t, "&")

			parseUrl := fmt.Sprintf("https://new.79da.com:665/m3u8.php?url=%s_%s",strings.Split(params[0],"=")[1],strings.Split(params[1],"=")[1] )


			return getRealUrl(parseUrl, stationInfo.Referer)
 		},
		 AfterFunc: func(stationInfo MovieStation, url string, resp *http.Response, w http.ResponseWriter, r *http.Request) {
			 defer resp.Body.Close()

			 body, err := ioutil.ReadAll(resp.Body)
			 if err != nil {
				 http.Error(w, err.Error(), 503)
				 return
			 }

			 prefix := ""
			 RootPath := false
			 if strings.HasPrefix(url, "https://tx.haihaiyu.com"){
				 prefix = "https://tx.haihaiyu.com"
				 RootPath = true
			 } else if strings.HasPrefix(url, "https://vod6.wenshibaowenbei.com") {
				 prefix = "https://vod6.wenshibaowenbei.com"
				 RootPath = true
			 }else {
				 path := url[0:strings.Index(url, "?")]
				 prefix = path[0:strings.LastIndex(url, "/")]
			 }
			 var newbody []byte
			 regexc := stationInfo.ReRegxp
			 if strings.HasSuffix(prefix, "/") && RootPath {
				 newbody = regexc.ReplaceAll(body, []byte(prefix[:len(prefix) - 1]+"$0"))
			 } else if strings.HasSuffix(prefix, "/") && !RootPath {
				 newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			 } else if !strings.HasSuffix(prefix, "/")  && RootPath {
				 newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			 } else {
				 newbody = regexc.ReplaceAll(body, []byte(prefix+"/$0"))
			 }
			 //	w.Header().Set("Content-Length", strconv.Itoa(len(newbody)))
			 w.Write(newbody)
		 },
	}
}

func getRealUrl(url string, referer string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Referer", referer)
	req.Header.Set("accept", `*/*`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return utils.MatchOneOf(string(body), `var vid="([^"]+)`)[1]
}



//https://www.jiujiukanpian.com/play/
type MovieStation struct {
	Name  string
	Prefix string
	Referer string
	DirectReturn bool
	HostUrl   string
	RootPath  bool
	ContentType string
	ReRegxp      *regexp.Regexp
	Jump    bool
	UrlBuildFunc func(requestUrl string, stationInfo MovieStation) string
	BeforeFunc func(stationInfo MovieStation, url string, header http.Header)
	AfterFunc  func(stationInfo MovieStation, url string, resp *http.Response, w http.ResponseWriter, r *http.Request)
}

var MoiveStationMap = make(map[string]*MovieStation, 0)

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	url := r.Form.Get("url")
	if url == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	hostPrefix := getHostPrefix(url)
	stationInfo := MoiveStationMap[hostPrefix]

	if stationInfo == nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	if stationInfo.UrlBuildFunc != nil {
		url = stationInfo.UrlBuildFunc(url, *stationInfo)
	}

	if stationInfo.DirectReturn {
		w.Header().Set("Content-Type", "audio/x-mpegurl")
		http.RedirectHandler(url, 302).ServeHTTP(w, r)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if stationInfo.BeforeFunc != nil {
		stationInfo.BeforeFunc(*stationInfo, url, req.Header)
	} else {
		if stationInfo.Referer != "" {
			req.Header.Set("Referer", stationInfo.Referer)
		}
		req.Header.Set("accept", `*/*`)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	if stationInfo.ContentType != "" {
		w.Header().Set("Content-Type", stationInfo.ContentType)
	} else {
		w.Header().Set("Content-Type", "audio/x-mpegurl")

	}

	if stationInfo.AfterFunc == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
		if stationInfo.Prefix != "" {
			prefix := stationInfo.Prefix
			var newbody []byte
			regexc := stationInfo.ReRegxp
			if strings.HasSuffix(prefix, "/") && stationInfo.RootPath {
				newbody = regexc.ReplaceAll(body, []byte(prefix[:len(prefix) - 1]+"$0"))
			} else if strings.HasSuffix(prefix, "/") && !stationInfo.RootPath {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			} else if !strings.HasSuffix(prefix, "/")  && stationInfo.RootPath {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			} else {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"/$0"))
			}
		//	w.Header().Set("Content-Length", strconv.Itoa(len(newbody)))
			w.Write(newbody)
		} else {
		//	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.Write(body)
		}
	} else {
		stationInfo.AfterFunc(*stationInfo, url, resp, w, r)
	}
}
	func getHostPrefix(url string) string {
		t := strings.Split(url, `//`)
		return t[0] + `//` + strings.Split(t[1], `/`)[0]
	}

