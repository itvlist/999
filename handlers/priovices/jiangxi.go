package priovices

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addJxtv("JXWS", "江西卫视 HD", "tv_jxtv1")
	addJxtv("JXDS", "江西都市频道 HD", "tv_jxtv2")
	addJxtv("JXJJSH", "江西经济生活频道 HD", "tv_jxtv3_hd")
	addJxtv("JXYSLY", "江西影视旅游频道 HD", "tv_jxtv4")
	addJxtv("JXGGNY", "江西公共农业频道 HD", "tv_jxtv5")
	addJxtv("JXSE", "江西少儿频道 HD", "tv_jxtv6")
	addJxtv("JXXW", "江西新闻频道 HD", "tv_jxtv7")
	addJxtv("JXYD", "江西移动电视 HD", "tv_jxtv8")
	addJxtv("JXFSGW", "江西风尚购物 HD", "tv_fsgw")
	addJxtv("JXTC", "江西陶瓷频道 HD", "tv_taoci")
}

type JxtvAuthResult struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	T      string `json:"t"`
	Count  int    `json:"count"`
	Token  string `json:"token"`
}

func addJxtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Referer: "http://www.jxntv.cn/",
		UrlFmt:  "https://live.jxtvcn.com.cn/live-jxtv/%s.m3u8?source=pc&t=%s&token=%s",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			charset := []rune("ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678oOLl9gqVvUuI1")
			eTags := make([]rune, 8)
			for i := 0; i < 8; i++ {
				eTags[i] = charset[rand.Intn(len(charset))]
			}
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			eTag := string(eTags)
			h := md5.New()
			h.Write([]byte(timestamp + id + ".m3u8" + eTag + "dfasg2df!f"))
			a := h.Sum(nil)

			values := url.Values{}
			values.Set("t", timestamp)
			values.Set("stream", id+".m3u8")

			req, err := http.NewRequest("POST", "https://app.jxntv.cn/Qiniu/liveauth/getPCAuth.php", strings.NewReader(values.Encode()))
			if err != nil {
				return ""
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			req.Header.Set("Origin", `http://www.jxntv.cn`)
			req.Header.Set("Referer", `http://www.jxntv.cn/`)
			req.Header.Set("Host", `app.jxntv.cn`)
			req.Header.Set("etag", eTag)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
			req.Header.Set("Content-Length", "34")
			req.Header.Set("Authorization", hex.EncodeToString(a))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return ""
			}

			result := JxtvAuthResult{}

			json.Unmarshal(body, &result)

			if result.Code != 200 {
				return ""
			}

			return fmt.Sprintf(refererInfo.UrlFmt, refererInfo.Id, result.T, result.Token)
		},
	}
}
