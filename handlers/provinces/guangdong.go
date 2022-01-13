package provinces

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addGdtv("GDWS", "广东卫视","43")
	addGdtv("GDZJ", "广东珠江","44")
	addGdtv("GDXW", "广东新闻","45")
	addGdtv("GDTY", "广东体育","47")
	addGdtv("NFWS", "广东南方卫视","51")
	addGdtv("GDJJKJ", "广东经济科教","49")
	addGdtv("GDYS", "广东影视","53")
	addGdtv("GDZY", "广东综艺","16")
	addGdtv("GDGJ", "广东珠江","46")
	addGdtv("GDSE", "广东少儿","54")
	addGdtv("GDJJKT", "广东嘉佳卡通","66")
	addGdtv("GDNFGW", "广东珠江","42")
	addGdtv("GDLNXQ", "广东岭南戏曲","15")
	addGdtv("GDFC", "广东房产","67")
	addGdtv("GDXDJY", "广东现代教育","13")
	addGdtv("GDYD", "广东移动","74")
	addGdtv("GRTNWHPD","广东GRTN文化频道", "75")
}

type nodeParamResult struct {
	Node string `json:"node"`
}

type gdtvQueryResult struct {
	AvatarUrl  string `json:"avatarUrl"`
	Category   int    `json:"category"`
	CoverUrl   string `json:"coverUrl"`
	Keyword    string `json:"keyword"`
	Name       string `json:"name"`
	Pk         int    `json:"pk"`
	PlayUrl    string `json:"playUrl"`
	Slogan     string `json:"slogan"`
	TimeOffset int    `json:"timeOffset"`
}
type gdtvPlayUrlResult struct {
	Hd string `json:"hd"`
}

func optionForUrl(url string) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("access-control-request-method", "GET")
	req.Header.Set("authority", "gdtv-api.gdtv.cn")
	req.Header.Set("access-control-request-headers", "content-type,x-itouchtv-ca-key,x-itouchtv-ca-signature,x-itouchtv-ca-timestamp,x-itouchtv-client,x-itouchtv-device-id")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("origin", `https://www.gdtv.cn`)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

}

//iv'1563432177954301' pkpadding7 004BC54E474B0F72902CE2E29B91C5E0A3E9E3BA11435370FAEAC96B25C9805D
func optionForGetParam(url string) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("access-control-request-method", "GET")
	req.Header.Set("accept", "*/*")
	req.Header.Set("authority", "tcdn-api.itouchtv.cn")
	req.Header.Set("access-control-request-headers", "content-type,x-itouchtv-ca-key,x-itouchtv-ca-signature,x-itouchtv-ca-timestamp,x-itouchtv-client,x-itouchtv-device-id")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("origin", `https://www.gdtv.cn`)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

}

func addGdtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Name: name,
		Referer: "http://app.cetv.cn",
		UrlFmt:  "https://gdtv-api.gdtv.cn/api/tv/v2/tvChannel/%s?tvChannelPk=%s&node=%s",
		Prefix:  "http://txycsbl.centv.cn/zb/",
		ReRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			optionForGetParam("https://tcdn-api.itouchtv.cn/getParam")
			req, err := http.NewRequest("GET", "https://tcdn-api.itouchtv.cn/getParam", nil)
			node := ""
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}

			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
			req.Header.Set("authority", "tcdn-api.itouchtv.cn")
			req.Header.Set("x-itouchtv-ca-timestamp", timestamp)
			req.Header.Set("x-itouchtv-ca-key", "89541443007807288657755311869534")
			req.Header.Set("x-itouchtv-client", "WEB_PC")
			req.Header.Set("x-itouchtv-device-id", "WEB_d547fdf0-633e-11ec-83d3-fb13b5511434")
			req.Header.Set("content-type", "application/json")
			req.Header.Set("accept", "application/json, text/plain, */*")
			req.Header.Set("origin", "https://www.gdtv.cn")
			req.Header.Set("sec-ch-ua-mobile", "?0")
			req.Header.Set("sec-ch-ua-platform", "macOS")
			req.Header.Set("sec-fetch-site", "same-site")
			req.Header.Set("sec-fetch-mode", "cors")
			req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
			secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
			message := []byte(fmt.Sprintf("GET\n%s\n%s\n", "https://tcdn-api.itouchtv.cn/getParam", timestamp))
			hash := hmac.New(sha256.New, secret)
			hash.Write(message)
			signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			req.Header.Set("x-itouchtv-ca-signature", signature)
			req.Header.Set("Referer", "https://www.gdtv.cn/")
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			paramResult := nodeParamResult{}
			err = json.Unmarshal(body, &paramResult)
			if err != nil || paramResult.Node == "" {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			} else {
				node = paramResult.Node
			}

			url := fmt.Sprintf(refererInfo.UrlFmt, id, id, base64.StdEncoding.EncodeToString([]byte(node)))
			optionForUrl(url)
			return url
		},
		BeforeFunc: func(refererInfo handlers.RerferInfo, url string, header http.Header) {
			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
			header.Set("authority", "gdtv-api.gdtv.cn")
			header.Set("x-itouchtv-ca-timestamp", timestamp)
			header.Set("x-itouchtv-ca-key", "89541443007807288657755311869534")
			header.Set("x-itouchtv-client", "WEB_PC")
			header.Set("x-itouchtv-device-id", "WEB_d547fdf0-633e-11ec-83d3-fb13b5511434")
			header.Set("content-type", "application/json")
			header.Set("accept", "application/json, text/plain, */*")
			header.Set("origin", "https://www.gdtv.cn")
			header.Set("sec-ch-ua-mobile", "?0")
			header.Set("sec-ch-ua-platform", "macOS")
			header.Set("sec-fetch-site", "same-site")
			header.Set("sec-fetch-mode", "cors")
			header.Set("accept-language", "zh-CN,zh;q=0.9")
			secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
			message := []byte(fmt.Sprintf("GET\n%s\n%s\n", url, timestamp))
			hash := hmac.New(sha256.New, secret)
			hash.Write(message)
			signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			header.Set("x-itouchtv-ca-signature", signature)
			header.Set("Referer", "https://www.gdtv.cn/")

		},
		AfterFunc: func(refererInfo handlers.RerferInfo, srcUrl string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}
			gdtvQueryResult := gdtvQueryResult{}
			err = json.Unmarshal(body, &gdtvQueryResult)
			gdtvPlayUrlResult := gdtvPlayUrlResult{}
			err = json.Unmarshal([]byte(gdtvQueryResult.PlayUrl), &gdtvPlayUrlResult)

			values := url.Values{}
			values.Set("url", gdtvPlayUrlResult.Hd)
			values.Set("referer", "https://www.gdtv.cn/")

			urlstr := fmt.Sprintf("http://%s:8880/transfer?%s", host, values.Encode())
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(urlstr, 302).ServeHTTP(w, r)
		},
	}
}
