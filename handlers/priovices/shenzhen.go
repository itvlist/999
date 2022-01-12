package priovices

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wmenjoy.com/iptv/handlers"
	"wmenjoy.com/iptv/utils"
)

func init()  {
	addSztv("SZWS", "深圳卫视", "AxeFRth")
	addSztv("SZYL", "深圳娱乐生活","1q4iPng")
	addSztv("SZSE", "深圳少儿","1SIQj6s")
	addSztv("SZGG", "深圳公共","2q76Sw2")
	addSztv("SZDSJ", "深圳电视剧","4azbkoY")
	addSztv("SZDB", "深圳未知","9zoW71b")
	addSztv("SZCJ", "深圳财经","3vlcoxP")
	addSztv("SZYHGW", "深圳宜和购物","BJ5u5k2")
	addSztv("SZDS", "深圳都市","ZwxzUXr")
	addSztv("SZGJ", "深圳国际","sztvgjpd")
	addSztv("SZTYJK","深圳体育教课", "sztvtyjk")
	addSztv("SZLG", "深圳未知","uGzbXhS")
	addSztv("SZYD", "深圳移动电视","wDF6KJ3")
	addSztv("SZDVSH", "深圳未知","xO1xQFv")
}

var reRegx = regexp.MustCompilePOSIX(`([^#]+\.ts)`)


func addSztv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Name: name,
		UrlFmt:  "https://sztv-live.cutv.com/%s/%s/%s.m3u8",
		ReRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			secret := "cutvLiveStream|Dream2017"
			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:10]
			h := md5.New()
			h.Write([]byte(timestamp + id + secret))
			a := h.Sum(nil)
			c := hex.EncodeToString(a)

			numUrl := fmt.Sprintf(
				"https://cls2.cutv.com/getCutvHlsLiveKey?t=%s&id=%s&token=%s&at=1", timestamp, id, c)
			req, err := http.NewRequest("GET", numUrl, nil)
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

			strBody := strings.Trim(string(body), "\"")
			fileName := ""
			index := len(strBody)
			if index > 0 {
				index -= index / 2
				strBody = strBody[index:] + strBody[0:index]
				fileNameb, _ := base64.StdEncoding.DecodeString(utils.ReverseString(strBody))
				fileName = string(fileNameb)
			}

			url := fmt.Sprintf(refererInfo.UrlFmt, id, "500", fileName)
			return url
		},
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()

			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}

			index := strings.LastIndex(url, "/")
			newbody := reRegx.ReplaceAll(body, []byte(url[0:index]+"/$0"))

			logrus.Info(string(newbody))
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			w.Write(newbody)
		},
	}
}
