package provinces

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addGztv("GZWS", "贵州卫视 HD", "ch01")
	addGztv("GZGG", "贵州公共 HD", "ch02")
	addGztv("GZYSWY", "贵州影视文艺 HD", "ch03")
	addGztv("GZDZSH", "贵州大众生活 HD", "ch04")
	addGztv("GZD5", "贵州第5频道 HD", "ch05")
	addGztv("GZKJJK", "贵州科教健康 HD", "ch06")
	addGztv("GZSZYD", "贵州数字移动 HD", "ch07")
}

type gztvResult struct {
	Title       string `json:"title"`
	EntryType   string `json:"entry_type"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	PubDate     string `json:"pub_date"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Author      string `json:"author"`
	StreamUrl   string `json:"stream_url"`
}

func addGztv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Referer: "https://www.gzstv.com/",
		UrlFmt:  "https://api.gzstv.com/v1/tv/%s/",
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := gztvResult{}

			json.Unmarshal(body, &result)

			http.RedirectHandler(result.StreamUrl, 302).ServeHTTP(w, r)
		},
	}
}
