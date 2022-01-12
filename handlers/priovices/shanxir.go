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
	addXiAntv("SXXAXWZH", "西安新闻综合", "xatv1", "综合", "720")
	addXiAntv("SXXADS", "西安都市频道", "xatv2", "综合", "720")
	addXiAntv("SXXASWZX", "西安商务资讯", "xatv3", "咨询", "720")
	addXiAntv("SXXAYS", "西安影视频道", "xatv4", "影视", "720")
	addXiAntv("SXXASL", "西安思路频道", "xatv5", "综合", "720")
	addXiAntv("SXXAYD", "西安移动频道", "xatv7", "综合", "720")
}

func addXiAntv(key string, name string, id string,category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		Group: "陕西",
		Category: category,
		Quality: quality,
		DirectReturn: true,
		UrlFmt:  "http://stream2.xiancity.cn/%s/hd/live.m3u8",
	}
}
