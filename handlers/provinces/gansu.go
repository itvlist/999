package provinces

import "wmenjoy.com/iptv/handlers"

func init(){
	//卫视,甘肃卫视 高清,http://39.134.39.39/PLTV/88888888/224/3221226240/index.m3u8
	//卫视,甘肃卫视 高清,http://223.110.243.215:80/PLTV/3/224/3221227568/index.m3u8
	addTvInfo("GSWS", 37, "甘肃卫视", "卫视","甘肃","综合","综合",
		"http://epg.51zmt.top:8000/tb1/ws/gansu.png")
	AddNanjingMobileTvChannel("GSWS", "3221227568", 50, "高清")
	AddLanzhouMobileTvChannel("GSWS", "3221226240", 60, "高清")
}

func AddNanjingMobileTvChannel(key string, id string, score int,  quality string) {
	_ = handlers.AddChannel(key, &handlers.IpTVChannel{
		Id:       id,
		Key:      key,
		Alive:	  true,
		Src:     handlers.ChannelSrcMap[handlers.NANJING_MOBILE_SRC],
		DirectReturn: true,
		Quality:  quality,
		Score:    score,
		IpList: []string{"223.110.243.215"},
		UrlFmt:   "/PLTV/3/224/%s/index.m3u8",
	})
}
func AddLanzhouMobileTvChannel(key string, id string, score int,  quality string) {
	_ = handlers.AddChannel(key, &handlers.IpTVChannel{
		Id:       id,
		Key:      key,
		Src:     handlers.ChannelSrcMap[handlers.LANZHOU_MOBILE_SRC],
		DirectReturn: true,
		Alive:	  true,
		Quality:  quality,
		Score:    score,
		IpList: []string{"39.134.39.39"},
		UrlFmt:   "/PLTV/88888888/224/%s/index.m3u8",
	})
}
