//
package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CDITVHandler(w http.ResponseWriter, r *http.Request) {
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

	url = fmt.Sprintf("https://www.cditv.cn/api.php?op=live&type=playTv&fluency=sd&videotype=m3u8&catid=192&id=%s", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("Accept", `*/*`)
	req.Header.Set("Referer", `https://www.cditv.cn/show-192-1-1.html`)
	req.Header.Set("X-Requested-With", `XMLHttpRequest`)
	req.Header.Set("sec-ch-ua-mobile", `?0`)
	req.Header.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36`)
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	/*url = fmt.Sprintf("http://hw-m-l.cztv.com/%s", r.URL.Path[4:])*/
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	realUrl := string(body)
	http.RedirectHandler(realUrl,302).ServeHTTP(w, r)

	//newbody := reRegx.ReplaceAll(body, []byte("/zj/channels/lantian/"+id+"/$0"))
	//	w.Header().Set("Host", "hw-m-l.cztv.com")
	//	w.Header().Set("Origin", "http://hw-m-l.cztv.com")
	//	w.Header().Set("Referer",  fmt.Sprintf("http://hw-m-l.cztv.com/channels/lantian/%s/", id))
//logrus.Info(string(newbody))
//	content := 	fmt.Sprintf("#EXTM3U\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=874000\n%s",realUrl);
//	w.Write([]byte(content))
}