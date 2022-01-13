package handlers

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var TvInfoMap = make(map[string]*IpTvInfo, 0)
var ChannelMap = make(map[string]IpTVChannelList, 0)

func AddChannel(key string, newChannel *IpTVChannel) error{
	tvInfo := TvInfoMap[key]
	if tvInfo == nil {
		return errors.New("频道信息不存在")
	}
	for _, channel := range tvInfo.Channels{
		if channel.Src == newChannel.Src && channel.Id == newChannel.Id {

			return errors.New("频道信息已经存在")
		}
	}
	newChannel.IpTvInfo = tvInfo

	tvInfo.Channels = append(tvInfo.Channels, newChannel)

	sort.Sort(tvInfo.Channels)
	return nil
}

type IpTVChannelList []*IpTVChannel
type IpTVChannel struct {
	Key string
	Src string
	Id  string
	Redirect bool
	Alive    bool
	IpTvInfo *IpTvInfo
	DirectReturn bool
	Quality  string
	IpList  []string
	Port 	int
	// 用于打分
	Score   int
	Proxy    string
	Protocol string
	Referer  string
	UrlFmt   string
	ReRegxp      *regexp.Regexp
	Prefix       string
	UrlBuildFunc func(channel IpTVChannel) string
	BeforeFunc func(channel IpTVChannel, url string, header http.Header)
	AfterFunc  func(channel IpTVChannel, url string, resp *http.Response, w http.ResponseWriter, r *http.Request)
}

func (s IpTVChannelList) Len() int           { return len(s) }
func (s IpTVChannelList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IpTVChannelList) Less(i, j int) bool { return s[i].Score > s[j].Score }

func(c IpTVChannel) getValidRequestUrl() string{

	if c.UrlBuildFunc != nil {
		return c.UrlBuildFunc(c)
	}

	path := ""
	if c.UrlFmt != "" {
		path = fmt.Sprintf(c.UrlFmt, c.Id)
	}

	if len(c.IpList) > 0 {

		port := c.Port
		if port <= 0 {
			port = 80
		}
		protocol := c.Protocol
		if protocol == "" {
			protocol = "http"
		}

		if !strings.HasPrefix(path,"/") {
			path = "/" + path
		}
		return fmt.Sprintf("%s://%s:%d%s",protocol,c.IpList[0], port, path)
	}

	return path
}

type IpTvInfo struct {
	TvTagId string
	Image string
	Key  string
	Name string
	Category string
	SubCategory string
	Group  string
	SubGroup string
	Channels IpTVChannelList
}

func (i IpTvInfo) getChannel() *IpTVChannel{
	if len(i.Channels) > 0 {
		return i.Channels[0]
	}
	return nil
}


func TVHandler(w http.ResponseWriter, r *http.Request) {
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

	tvInfo := TvInfoMap[strings.ToUpper(id)]

	if tvInfo == nil  {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	tvChannel := tvInfo.getChannel()
	if tvChannel == nil  {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	url = tvChannel.getValidRequestUrl()

	if tvChannel.DirectReturn {
		w.Header().Set("Content-Type", "audio/x-mpegurl")
		http.RedirectHandler(url, 302).ServeHTTP(w, r)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if tvChannel.BeforeFunc != nil {
		tvChannel.BeforeFunc(*tvChannel, url, req.Header)
	} else {
		if tvChannel.Referer != "" {
			req.Header.Set("Referer", tvChannel.Referer)
		}
		req.Header.Set("accept", `*/*`)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	if tvChannel.AfterFunc == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
		if tvChannel.Prefix != "" {
			prefix := tvChannel.Prefix

			if strings.Contains(prefix, "%s") {
				prefix = fmt.Sprintf(prefix, tvChannel.Id)
			}

			w.Header().Set("Content-Type", "audio/x-mpegurl")
			var newbody []byte
			regexc := tvChannel.ReRegxp
			if regexc == nil {
				regexc = reRegx
			}

			if strings.HasSuffix(prefix, "/") {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			} else {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"/$0"))
			}
			w.Write(newbody)
		} else {
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			w.Write(body)
		}
	} else {
		tvChannel.AfterFunc(*tvChannel, url, resp, w, r)
	}

}