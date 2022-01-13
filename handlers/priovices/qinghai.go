package priovices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addQhtv("QHWS", "青海卫视","3", "卫视","标清")
	addQhtv("QHJS", "青海经视","2", "财经","标清")
	addQhtv("QHDS", "青海都市","4", "综合","标清")
	addQhtv("QHXWZH", "青海新闻综合","20", "新闻","标清")

}
type QinghaiTvInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Logo struct {
		Rectangle struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"rectangle"`
		Square struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"square"`
	} `json:"logo"`
	Snap struct {
		Host     string `json:"host"`
		Dir      string `json:"dir"`
		Filepath string `json:"filepath"`
		Filename string `json:"filename"`
	} `json:"snap"`
	M3U8       string `json:"m3u8"`
	CurProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"cur_program"`
	SaveTime    string `json:"save_time"`
	NextProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"next_program"`
	AudioOnly     string `json:"audio_only"`
	Aspect        string `json:"aspect"`
	ContentUrl    string `json:"content_url"`
	ChannelStream []struct {
		Url        string `json:"url"`
		Name       string `json:"name"`
		StreamName string `json:"stream_name"`
		M3U8       string `json:"m3u8"`
		Bitrate    string `json:"bitrate"`
	} `json:"channel_stream"`
	ShareUrl interface{} `json:"share_url"`
	Title    string      `json:"title"`
}
func addQhtv(key string, name string, id string,category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Group: "青海",
		Category: category,
		Quality: quality,
		DirectReturn: true,
		UrlFmt:  "http://www.qhbtv.com/m2o/channel/channel_info.php?id=%s",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			url := fmt.Sprintf(refererInfo.UrlFmt, id)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return ""
			}
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			req.Header.Set("Host", "http://www.qhbtv.com")
			req.Header.Set("access-control-request-headers", "content-type,x-itouchtv-ca-key,x-itouchtv-ca-signature,x-itouchtv-ca-timestamp,x-itouchtv-client,x-itouchtv-device-id")
			req.Header.Set("Referer", `http://www.qhbtv.com/`)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}
			result := make([]QinghaiTvInfo,0)

			json.Unmarshal(body, &result)

			return result[0].M3U8

		},
	}
}