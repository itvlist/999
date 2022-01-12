package priovices

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addAhtv("AHWS", "安徽卫视 HD", "47")
	addAhtv("AHJJSH", "安徽经济生活 HD", "71")
	addAhtv("AHZYTY", "安徽综艺体育 HD", "73")
	addAhtv("AHYS", "安徽影视 HD", "72")
	addAhtv("AHGJ", "安徽国际 HD", "50")
	addAhtv("AHNYKJ", "安徽农业科教 HD", "51")
	addAhtv("AHGJ", "安徽国际 HD", "70")
	addAhtv("AHYD", "安徽移动电视 HD", "68")
	addAhtv("AHJC", "睛彩安徽 HD", "85")
}

type anhuiTvResult struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Aspect     string `json:"aspect"`
	CommentNum int    `json:"comment_num"`
	Snap       struct {
		Host     string `json:"host"`
		Dir      string `json:"dir"`
		Path     string `json:"path"`
		Filepath string `json:"filepath"`
		Filename string `json:"filename"`
	} `json:"snap"`
	SiteId     int    `json:"site_id"`
	ClickNum   int    `json:"click_num"`
	PraiseNum  int    `json:"praise_num"`
	ShareNum   int    `json:"share_num"`
	M3U8       string `json:"m3u8"`
	CurProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"cur_program"`
	NextProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"next_program"`
	Logo struct {
		Square struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Path     string `json:"path"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"square"`
		Rectangle struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Path     string `json:"path"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"rectangle"`
	} `json:"logo"`
	SaveTime      string `json:"save_time"`
	AudioOnly     string `json:"audio_only"`
	ContentUrl    string `json:"content_url"`
	ChannelStream []struct {
		Url        string `json:"url"`
		Name       string `json:"name"`
		StreamName string `json:"stream_name"`
		M3U8       string `json:"m3u8"`
		Bitrate    string `json:"bitrate"`
		StreamUrl  string `json:"stream_url"`
	} `json:"channel_stream"`
	NodeId   int    `json:"node_id"`
	ShareUrl string `json:"share_url"`
	OriginId int    `json:"origin_id"`
}

//        "appid":"m2otdjzyuuu8bcccnq",
//                    "appkey":"5eab6b4e1969a8f9aef459699f0d9000",
func addAhtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Referer: "http://www.ahtv.cn/",
		UrlFmt:  "http://mapi.ahtv.cn/api/open/ahtv/channel.php?appid=m2otdjzyuuu8bcccnq&appkey=5eab6b4e1969a8f9aef459699f0d9000&is_audio=0&category_id=1%2C2",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			return refererInfo.UrlFmt
		},
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := make([]anhuiTvResult, 0)

			_ = json.Unmarshal(body, &result)
			id, _ := strconv.Atoi(refererInfo.Id)
			for _, value := range result {
				if value.Id == id {
					w.Header().Set("Content-Type", "audio/x-mpegurl")
					http.RedirectHandler(string(value.M3U8), 302).ServeHTTP(w, r)
					return
				}
			}
		},
	}
}
