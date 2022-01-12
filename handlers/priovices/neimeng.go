package priovices

import (
	"fmt"
	"strings"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addNmtv("NMWS", "内蒙古卫视 HD", "live9:nmgws")
	addNmtv("NMSE", "内蒙古少儿频道 HD", "live11:sepd")
	addNmtv("NMJJSH", "内蒙古经济生活 HD", "live10:jjsh")
	addNmtv("NMMGYWS", "内蒙古蒙古语卫视 HD", "live9:nmgmgyws")
	addNmtv("NMMGYWH", "内蒙古蒙古语文化 HD", "live9:nmgmgywh")
	addNmtv("NMNM", "内蒙古农牧频道 HD", "live11:nmpd")
	addNmtv("NMXWZH", "内蒙古新闻综合 HD", "live10:xwzh")
	addNmtv("NMWTYL", "内蒙古文体娱乐 HD", "live10:wtyl")
}

/**
内蒙古卫视,http://live9.m2oplus.nmtv.cn/nmgws/playlist.m3u8
内蒙古卫视,http://live9.m2oplus.nmtv.cn/nmgws/hd/live.m3u8
内蒙古少儿频道,http://live11.m2oplus.nmtv.cn/sepd/hd/live.m3u8
内蒙古经济生活,http://live10.m2oplus.nmtv.cn/jjsh/hd/live.m3u8
内蒙古蒙古语卫视,http://live9.m2oplus.nmtv.cn/nmgmgyws/hd/live.m3u8
内蒙古蒙古语文化,http://live9.m2oplus.nmtv.cn/nmgmgywh/hd/live.m3u8
内蒙古农牧频道,http://live11.m2oplus.nmtv.cn/nmpd/hd/live.m3u8
内蒙古新闻综合,http://live10.m2oplus.nmtv.cn/xwzh/hd/live.m3u8
内蒙古文体娱乐,http://live10.m2oplus.nmtv.cn/wtyl/hd/live.m3u8
 */
func addNmtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		DirectReturn: true,
		UrlFmt:  "http://%s.m2oplus.nmtv.cn/%s/hd/live.m3u8",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			t := strings.Split(id, `:`)
			return fmt.Sprintf(refererInfo.UrlFmt,t[0],t[1])
		},
	}
}
