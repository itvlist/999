package provinces

import "wmenjoy.com/iptv/handlers"

func init(){
	//卫视,甘肃卫视 高清,http://39.134.39.39/PLTV/88888888/224/3221226240/index.m3u8
	//卫视,甘肃卫视 高清,http://223.110.243.215:80/PLTV/3/224/3221227568/index.m3u8
	addTvInfo("GSWS", "甘肃卫视", "卫视","甘肃","综合","综合")
	AddNanjingMobileTvChannel("GSWS", "3221227568", 50, "高清")
	AddLanzhouMobileTvChannel("GSWS", "3221226240", 60, "高清")
}

func AddNanjingMobileTvChannel(key string, id string, score int,  quality string) {
	_ = handlers.AddChannel(key, &handlers.IpTVChannel{
		Id:       id,
		Key:      key,
		Alive:	  true,
		Src:      "NanJingMobile",
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
		Src:      "LanZhouMobile",
		DirectReturn: true,
		Alive:	  true,
		Quality:  quality,
		Score:    score,
		IpList: []string{"39.134.39.39"},
		UrlFmt:   "/PLTV/88888888/224/%s/index.m3u8",
	})
}
