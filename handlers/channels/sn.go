package channels

import (
	"wmenjoy.com/iptv/handlers"
)

func addSnChannel(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:           id,
		Name:         name,
		Group:        "央视",
		Category:     category,
		Quality:      quality,
		UrlFmt:       "http://dbiptv.sn.chinamobile.com/PLTV/88888888/224/%s/1.m3u8",
		DirectReturn: false,
	}
}

func init() {
	addSnChannel("CHC-DZ", "CHC动作电影", "3221226465", "影视", "超清")
}
