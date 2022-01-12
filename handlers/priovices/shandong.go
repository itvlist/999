package priovices

import "wmenjoy.com/iptv/handlers"

/**
山东,山东卫视,https://livealone302.iqilu.com/iqilu/sdtv.m3u8
山东,山东齐鲁频道,https://livealone302.iqilu.com/iqilu/qlpd.m3u8
山东,山东体育频道,https://livealone302.iqilu.com/iqilu/typd.m3u8
山东,山东生活频道,https://livealone302.iqilu.com/iqilu/shpd.m3u8
山东,山东综艺频道,https://livealone302.iqilu.com/iqilu/zypd.m3u8
山东,山东新闻频道,https://livealone302.iqilu.com/iqilu/ggpd.m3u8
山东,山东少儿频道,https://livealone302.iqilu.com/iqilu/sepd.m3u8
山东,山东文旅频道,https://livealone302.iqilu.com/iqilu/yspd.m3u8
山东,山东农科频道,https://livealone302.iqilu.com/iqilu/nkpd.m3u8
山东,山东体育频道,https://livealone302.iqilu.com/iqilu/typd.m3u8
*/
func init()  {
	addSdtv("SDWS", "山东卫视", "sdtv")
	addSdtv("SDQL", "山东齐鲁频道", "qlpd")
	addSdtv("SDTY", "山东体育频道", "typd")
	addSdtv("SDSH", "山东生活频道", "shpd")
	addSdtv("SDZY", "山东综艺频道", "zypd")
	addSdtv("SDXW", "山东新闻频道", "ggpd")
	addSdtv("SDSE", "山东少儿频道", "sepd")
	addSdtv("SDWL", "山东文旅频道", "yspd")
	addSdtv("SDNK", "山东农科频道", "nkpd")
	addSdtv("SDTY", "山东体育频道", "typd")
}

func addSdtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		DirectReturn: true,
		UrlFmt:  "https://livealone302.iqilu.com/iqilu/%s.m3u8",
	}
}

