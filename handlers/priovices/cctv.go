package priovices

import (
	"regexp"
	"wmenjoy.com/iptv/handlers"
)
// https://cctvalih5c.v.myalicdn.com/live/cdrmcctv17_1td.m3u8

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


func addCctv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Name:	name,
		Group: "央视",
		Category: category,
		Quality: quality,
		UrlFmt:  "https://cctvalih5c.v.myalicdn.com/live/cdrm%s_1td.m3u8",
		DirectReturn: true,
	}
}

func init() {
	addCetv("CETV-1", "中国教育电视台1", "cetv1")
	addCetv("CETV-2", "中国教育电视台1", "cetv2")
	addCetv("CETV-3", "中国教育电视台1", "cetv3")
	addCetv("CETV-4", "中国教育电视台1", "cetv4")

	addCctv("CCTV-1", "CCTV-1 综合", "cctv1", "综合","超清")
	addCctv("CCTV-2", "CCTV-2 财经", "cctv2", "财经","超清")
	addCctv("CCTV-3", "CCTV-3 综艺", "cctv3", "综艺","超清")
	addCctv("CCTV-4", "CCTV-4 中文国际", "cctv4", "国际","超清")
	addCctv("CCTV-5", "CCTV-5 体育", "cctv5", "体育","超清")
	addCctv("CCTV-5+", "CCTV-5+ 体育", "cctv5plus", "体育","超清")
	addCctv("CCTV-6", "CCTV-6 电影", "cctv6", "影视","超清")
	addCctv("CCTV-7", "CCTV-7 国防军事", "cctv7", "军事","超清")
	addCctv("CCTV-8", "CCTV-8 电视剧", "cctv8", "影视","超清")
	addCctv("CCTV-9", "CCTV-9 记录", "cctv9", "科教","超清")
	addCctv("CCTV-10", "CCTV-10 科教", "cctv10", "科教","超清")
	addCctv("CCTV-11", "CCTV-11 戏曲", "cctv11", "综艺","超清")
	addCctv("CCTV-12", "CCTV-12 社会与法", "cctv12", "法律","超清")
	addCctv("CCTV-13", "CCTV-13 新闻", "cctv13", "新闻","超清")
	addCctv("CCTV-14", "CCTV-14 少儿", "cctv14", "少儿","超清")
	addCctv("CCTV-15", "CCTV-15 音乐", "cctv15", "综艺","超清")
	addCctv("CCTV-16", "CCTV-16 奥林匹克", "cctv16", "体育","超清")
	addCctv("CCTV-17", "CCTV-17 农业农村", "cctv17", "农村","超清")
	addCctv("CCTV-4A", "CCTV-4A 美国", "cctvamerica", "国际","超清")
	addCctv("CCTV-4E", "CCTV-4E 欧洲", "cctv1europe", "国际","超清")

}
