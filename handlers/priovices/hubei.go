package priovices

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addHbtv("HBLS", "湖北垄上 HD", "438")
	addHbtv("MJGW", "湖北美嘉购物 HD", "439")
	addHbtv("HBJY", "湖北教育 HD", "437")
	addHbtv("HBSH", "湖北生活 HD", "436")
	addHbtv("HBYS", "湖北影视 HD", "435")
	addHbtv("HBGG", "湖北公共 HD", "434")
	addHbtv("HBZH", "湖北综合 HD", "433")
	addHbtv("HBJS", "湖北经视 HD", "432")
	addHbtv("HBWS", "湖北卫视 HD", "431")
	addHbtv("HBXW", "湖北新闻 HD", "470")
}

type HbTvResult struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Alias          string `json:"alias"`
	Stream         string `json:"stream"`
	Icon           string `json:"icon"`
	LiveType       string `json:"live_type"`
	AccessType     string `json:"access_type"`
	Fms            string `json:"fms"`
	Created        string `json:"created"`
	Createdby      string `json:"createdby"`
	PlaybillStatus string `json:"playbill_status"`
	Replay         string `json:"replay"`
	ReplayExpire   string `json:"replay_expire"`
	HasHd          string `json:"has_hd"`
	PublishedPc    string `json:"published_pc"`
	PublishedPhone string `json:"published_phone"`
	Rate           string `json:"rate"`
	DefaultThumb   string `json:"default_thumb"`
	StreamHd       string `json:"stream_hd"`
	Rtmp           string `json:"rtmp"`
	RtmpHd         string `json:"rtmp_hd"`
	Sort           string `json:"sort"`
	State          string `json:"state"`
	Url            string `json:"url"`
	PlayUrl        string `json:"play_url"`
	PlayUrlSd      string `json:"play_url_sd"`
	PlayUrlHd      string `json:"play_url_hd"`
	UsePub         string `json:"use_pub"`
	Shift          string `json:"shift"`
	ShiftStarttime string `json:"shift_starttime"`
	ShiftEndtime   string `json:"shift_endtime"`
	VirtualLive    string `json:"virtual_live"`
	PlayControl    string `json:"play_control"`
	VmsTid         string `json:"vms_tid"`
	ControlUrl     string `json:"control_url"`
	ShiftDay       string `json:"shift_day"`
	Token          string `json:"token"`
	Identity       string `json:"identity"`
	OssId          string `json:"oss_id"`
	OssHdId        string `json:"oss_hd_id"`
	Type           string `json:"type"`
}

func addHbtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Referer: "http://app.cjyun.org/",
		UrlFmt:  "http://app.cjyun.org/video/player/stream?stream_id=%s&site_id=10008",
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := HbTvResult{}

			_ = json.Unmarshal(body, &result)

			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(string(result.Stream), 302).ServeHTTP(w, r)
		},
	}
}

