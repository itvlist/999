package provinces

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"wmenjoy.com/iptv/handlers"
	"wmenjoy.com/iptv/utils"
)

func init()  {
	//加密了，不在使用
	addCqtv("CQWS", "重庆卫视 HD", "4918")
}
func getKey() string {
	url := "https://sjlivecdnx.cbg.cn/1ive/stream_2.php"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("referer", "https://www.cbg.cn/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(body)
}

func addCqtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:     id,
		Key:    key,
		Name:   name,
		UrlFmt: "https://web.cbg.cn/live/getLiveUrl?url=%s",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {

			url := fmt.Sprintf("https://rmtapi.cbg.cn/list/%s/1.html?pagesize=20", refererInfo.Id)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return ""
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}
			urlResult := utils.MatchOneOf(string(body), "\"ios_HDlive_url\":\"([^\"]+)")[1]
			return fmt.Sprintf(refererInfo.UrlFmt, urlResult)
		},
		AfterFunc: func(refererInfo handlers.RerferInfo, srcUrl string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}
			urlResult := utils.MatchOneOf(string(body), "\"url\"[ ]*:[ ]*\"([^\"]+)")[1]
			proxyUrl := fmt.Sprintf("http://%s:8880/ats?file=", host)
			index := strings.LastIndex(urlResult, "/")
			values := url.Values{}
			values.Set("url", urlResult)
			values.Set("srcUrl", urlResult[0:index])
			values.Set("proxy", proxyUrl)
			values.Set("key", getKey())
			values.Set("referer", "https://www.cbg.cn/")
			realUrl := fmt.Sprintf("http://%s:8880/transfer?%s", host, values.Encode())
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(realUrl, 302).ServeHTTP(w, r)
		},
	}
}
