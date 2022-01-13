package provinces

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"wmenjoy.com/iptv/handlers"
)

func init()  {
	addJsws("JSWS", "江苏卫视 HD", "jswspindao")
	addJsws("JSCS", "江苏城市 HD", "jscspindao")
	addJsws("JSGG", "江苏公共 HD", "jsggpindao")
	addJsws("JSZY", "江苏综艺 HD", "jszypindao")
	addJsws("JSYS", "江苏影视 HD", "jsyspindao")
	addJsws("JSTY", "江苏体育休闲 HD", "jstypindao")
	addJsws("JSJY", "江苏教育 HD", "jsjypindao")
	addJsws("JSXX", "江苏学习 HD", "jsxxpindao")
	addJsws("JSGJ", "江苏国际 HD", "jsgjpindao")
	addJsws("JSLZ", "江苏靓装 HD", "jslzpindao")
	addJsws("HXGW", "好享购物 HD", "hxgwpindao")
	addJsws("YMKT", "优漫卡通 HD", "ymktpindao")
}
func addJsws(key string, name string, id string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Name:    name,
		DirectReturn: true,
		Referer: "http://live.jstv.com/",
		UrlFmt:  "https://live2021hls-yf.jstv.com/live/%s-1000/online.m3u8",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			e := fmt.Sprintf(refererInfo.UrlFmt, id)
			t := strings.Split(e, `//`)
			n := t[0] + `//` + strings.Split(t[1], `/`)[0]
			a := strings.ReplaceAll(e, n,"")
			r := "jstv123qwe"
			i := time.Now().Unix() + 300
			s := fmt.Sprintf("%s&%d&%s",r,i,a)
			h := md5.New()
			h.Write([]byte(s))
			md5x := h.Sum(nil)
			c :=hex.EncodeToString(md5x)
			return fmt.Sprintf("%s?upt=%s%d",e,c[12:20],i)

		},
	}
}
