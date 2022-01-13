package provinces

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
	addJxtv("JXWS", "江西卫视 HD", "tv_jxtv1","卫视","720")
	addJxtv("JXDS", "江西都市频道 HD", "tv_jxtv2","综合","720")
	addJxtv("JXJJSH", "江西经济生活频道 HD", "tv_jxtv3_hd","财经","720")
	addJxtv("JXYSLY", "江西影视旅游频道 HD", "tv_jxtv4","文旅","720")
	addJxtv("JXGGNY", "江西公共农业频道 HD", "tv_jxtv5","农业","720")
	addJxtv("JXSE", "江西少儿频道 HD", "tv_jxtv6","少儿","720")
	addJxtv("JXXW", "江西新闻频道 HD", "tv_jxtv7","新闻","720")
	addJxtv("JXYD", "江西移动电视 HD", "tv_jxtv8","综合","720")
	addJxtv("JXFSGW", "江西风尚购物 HD", "tv_fsgw","购物","720")
	addJxtv("JXTC", "江西陶瓷频道 HD", "tv_taoci","其他","720")
}

type JxtvAuthResult struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	T      string `json:"t"`
	Count  int    `json:"count"`
	Token  string `json:"token"`
}

func addJxtv(key string, name string, id string,category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Group: "江西",
		Category: category,
		Quality: quality,
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
