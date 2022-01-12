package priovices

import "wmenjoy.com/iptv/handlers"

/**
天津少儿,http://60.26.15.147:9901/tsfile/live/0110_1.m3u8
天津体育,http://60.26.15.147:9901/tsfile/live/0111_1.m3u8
天津文艺,http://60.26.15.147:9901/tsfile/live/0119_1.m3u8
天津新闻,http://60.26.15.147:9901/tsfile/live/0120_1.m3u8
天津卫视,http://60.26.15.147:9901/tsfile/live/0135_1.m3u8
天津影视,http://60.26.15.147:9901/tsfile/live/0136_1.m3u8
天津都市,http://60.26.15.147:9901/tsfile/live/0141_1.m3u8

 */
func init(){
	addTjtv("TJWS", "天津卫视", "0135", "卫视", "高清")
	addTjtv("TJWY", "天津文艺", "0119", "综艺", "高清")
	addTjtv("TJXW", "天津新闻", "0120", "新闻", "高清")
	addTjtv("TJDS", "天津都市", "0141", "综合", "高清")
	addTjtv("TJYS", "天津影视", "0136", "影视", "高清")
	addTjtv("TJTY", "天津体育", "0111", "体育", "高清")
	addTjtv("TJSE", "天津少儿", "0110", "少儿", "高清")
}

func addTjtv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Group: "天津",
		Category: category,
		Quality: quality,
		Name:    name,
		DirectReturn: true,
		UrlFmt:  "http://60.26.15.147:9901/tsfile/live/%s_1.m3u8",
	}
}

//天津影视HD,http://60.26.15.147:9901/tsfile/live/0136_1.m3u8