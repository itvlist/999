package handlers

import (
	"io"
	"net/http"
	"regexp"
)

type IpTVChannel struct {
	Src string
	Id  string
	Redirect bool
	DirectReturn bool
	Quality  string
	IpList  []string
	Path 	string
	// 用于打分
	Score   int
	Proxy    string
	Referer  string
	UrlFmt   string
	ReRegxp      *regexp.Regexp
	Prefix       string
	UrlBuildFunc func(channel IpTVChannel) string
	BeforeFunc func(channel IpTVChannel, url string, header http.Header)
	AfterFunc  func(channel IpTVChannel, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request)
}

type IpTvInfo struct {
	Key  string
	Name string
	Category string
	SubCategory string
	Group  string
	SubGroup string
	Channels [] IpTVChannel
}
