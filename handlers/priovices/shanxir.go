package priovices

import (
	"wmenjoy.com/iptv/handlers"
)
/**
陕西,西安新闻综合,http://stream2.xiancity.cn/xatv1/hd/live.m3u8
陕西,西安都市频道,http://stream2.xiancity.cn/xatv2/sd/live.m3u8
陕西,西安商务资讯,http://stream2.xiancity.cn/xatv3/sd/live.m3u8
陕西,西安影视频道,http://stream2.xiancity.cn/xatv4/sd/live.m3u8
陕西,西安思路频道,http://stream2.xiancity.cn/xatv5/sd/live.m3u8
陕西,西安移动频道,http://stream2.xiancity.cn/xatv7/sd/live.m3u8
 */
func init()  {
	addXiAntv("SXXAXWZH", "西安新闻综合", "xatv1")
	addXiAntv("SXXADS", "西安都市频道", "xatv2")
	addXiAntv("SXXASWZX", "西安商务资讯", "xatv3")
	addXiAntv("SXXAYS", "西安影视频道", "xatv4")
	addXiAntv("SXXASL", "西安思路频道", "xatv5")
	addXiAntv("SXXAYD", "西安移动频道", "xatv7")
}

func addXiAntv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		DirectReturn: true,
		UrlFmt:  "http://stream2.xiancity.cn/%s/hd/live.m3u8",
	}
}
