package provinces

import (
	"fmt"
	"strings"
	"wmenjoy.com/iptv/handlers"
)
/**
河北,河北卫视,http://live2.plus.hebtv.com/hbwsx/hd/live.m3u8
河北,河北经济生活,http://live2.plus.hebtv.com/jjshx/hd/live.m3u8
河北,河北农民频道,http://live3.plus.hebtv.com/nmpdx/hd/live.m3u8
河北,河北都市,http://live3.plus.hebtv.com/hbdsx/hd/live.m3u8
河北,河北影视,http://live6.plus.hebtv.com/hbysx/hd/live.m3u8
河北,河北少儿科教,http://live6.plus.hebtv.com/sekjx/hd/live.m3u8
河北,河北公共频道,http://live7.plus.hebtv.com/hbggx/hd/live.m3u8
*/
func init()  {
	addHebtv("HEBWS", "河北卫视 720", "live2:hbwsx")
	addHebtv("HEBJJSH", "河北经济生活 720", "live2:jjshx")
	addHebtv("HEBNM", "河北农民频道 720", "live3:nmpdx")
	addHebtv("HEBDS", "河北都市 720", "live3:hbdsx")
	addHebtv("HEBYS", "河北影视 720", "live6:hbysx")
	addHebtv("HEBSEKJ", "河北少儿科教 720", "live6:sekjx")
	addHebtv("HEBGG", "河北公共频道 720", "live7:hbggx")

}

func addHebtv(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		DirectReturn: true,
		UrlFmt:  "http://%s.plus.hebtv.com/%s/hd/live.m3u8",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			t := strings.Split(id, `:`)
			return fmt.Sprintf(refererInfo.UrlFmt,t[0],t[1])
		},
	}
}

