package priovices

import (
	"regexp"
	"wmenjoy.com/iptv/handlers"
)

func addCetv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Name:	name,
		Referer: "http://app.cetv.cn",
		UrlFmt:  "http://txycsbl.centv.cn/zb/0104%s.m3u8",
		Prefix:  "http://txycsbl.centv.cn/zb/",
		ReRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
	}
}

func init() {
	addCetv("CETV-1", "中国教育电视台1", "cetv1")
	addCetv("CETV-2", "中国教育电视台1", "cetv2")
	addCetv("CETV-3", "中国教育电视台1", "cetv3")
	addCetv("CETV-4", "中国教育电视台1", "cetv4")
}
