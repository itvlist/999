package provinces

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

// 加密了，目前需要考虑其他办法
func addCctv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Name:	name,
		Group: "央视",
		Category: category,
		Quality: quality,
		UrlFmt:  "https://cctvtxyh5ca.liveplay.myqcloud.com/live/%s/index.m3u8 ",
		DirectReturn: true,
	}
}

func init() {
	addCetv("CETV-1", "中国教育电视台1", "cetv1")
	addCetv("CETV-2", "中国教育电视台1", "cetv2")
	addCetv("CETV-3", "中国教育电视台1", "cetv3")
	addCetv("CETV-4", "中国教育电视台1", "cetv4")

	addCctv("CCTV-1", "CCTV-1 综合", "cctv1_2_hd", "综合","超清")
	addCctv("CCTV-2", "CCTV-2 财经", "cctv2_2_hd", "财经","超清")
	addCctv("CCTV-3", "CCTV-3 综艺", "cctv3_2_hd", "综艺","超清")
	addCctv("CCTV-4", "CCTV-4 中文国际", "cctv4_2_hd", "国际","超清")
	addCctv("CCTV-5", "CCTV-5 体育", "cctv5_2", "体育","超清")
	addCctv("CCTV-5+", "CCTV-5+ 体育", "cctv5plus_2", "体育","超清")
	addCctv("CCTV-6", "CCTV-6 电影", "cctv6_2_hd", "影视","超清")
	addCctv("CCTV-7", "CCTV-7 国防军事", "cctv7_2_hd", "军事","超清")
	addCctv("CCTV-8", "CCTV-8 电视剧", "cctv8_2_hd", "影视","超清")
	addCctv("CCTV-9", "CCTV-9 记录", "cdrmcctvjilu_1", "科教","超清")
	addCctv("CCTV-10", "CCTV-10 科教", "cctv10_2_hd", "科教","超清")
	addCctv("CCTV-11", "CCTV-11 戏曲", "cctv11_2_hd", "综艺","超清")
	addCctv("CCTV-12", "CCTV-12 社会与法", "cctv12_2_hd", "法律","超清")
	addCctv("CCTV-13", "CCTV-13 新闻", "cctv13_2_hd", "新闻","超清")
	addCctv("CCTV-14", "CCTV-14 少儿", "cdrmcctvchild_1", "少儿","超清")
	addCctv("CCTV-15", "CCTV-15 音乐", "cctv15_2_hd", "综艺","超清")
	addCctv("CCTV-16", "CCTV-16 奥林匹克", "cctv16_2", "体育","超清")
	addCctv("CCTV-17", "CCTV-17 农业农村", "cctv17_2_hd", "农村","超清")
	addCctv("CCTV-4A", "CCTV-4A 美国", "cdrmcctvamerica_1", "国际","超清")
	addCctv("CCTV-4E", "CCTV-4E 欧洲", "cdrmcctveurope_1", "国际","超清")

}
