//
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var idmapper = map[string]string{
	"257":"1500",
}

type (
	HTTPResult struct {
		Code  int `json:"code"`
		ErrCode int `json:"errCode"`
		Data ResultData
		ErrMessage string `json:"errMessage"`
	}

	ResultData struct {
		HlsUrl string `json:"hlsurl"`
	}
)

func HangZhouhander(w http.ResponseWriter, r *http.Request) {
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

	url = fmt.Sprintf("https://1812501212048408.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/node-api.online/node-api/tv/getPlayAddress?id=%s&uin=%s&clientId=b6c4e124-75e7-4cba-be86-31d45a73472f", id, idmapper[id])
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("Accept", `application/json, text/plain, */*`)
	req.Header.Set("Referer", `https://web.guangdianyun.tv/`)
	req.Header.Set("X-Ca-Stage", ``)
	req.Header.Set("token", ``)
	req.Header.Set("Connection", `keep-alive`)
	req.Header.Set("Origin", `https://web.guangdianyun.tv`)

	req.Header.Set("Sec-Fetch-Site", `cross-site`)
	req.Header.Set("Sec-Fetch-Mode", `cors`)
	req.Header.Set("Sec-Fetch-Dest", `empty`)
	req.Header.Set("Accept-Language", `zh-CN,zh;q=0.9`)
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
	var mapResult HTTPResult
	err = json.Unmarshal(body, &mapResult)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	http.RedirectHandler(mapResult.Data.HlsUrl,302).ServeHTTP(w, r)

	//newbody := reRegx.ReplaceAll(body, []byte("/zj/channels/lantian/"+id+"/$0"))
	//	w.Header().Set("Host", "hw-m-l.cztv.com")
	//	w.Header().Set("Origin", "http://hw-m-l.cztv.com")
	//	w.Header().Set("Referer",  fmt.Sprintf("http://hw-m-l.cztv.com/channels/lantian/%s/", id))
	//logrus.Info(string(newbody))
	//	content := 	fmt.Sprintf("#EXTM3U\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=874000\n%s",realUrl);
	//	w.Write([]byte(content))
}