package priovices

import (
	"io"
	"io/ioutil"
	"net/http"
	"wmenjoy.com/iptv/handlers"
)
/**
四川妇女儿童,http://scgctvshow.sctv.com/hdlive/sctv7/1.m3u8
四川妇女儿童,http://scgctvshow.sctv.com/hdlive/sctv7/index.m3u8
四川影视文化,http://scgctvshow.sctv.com/hdlive/sctv5/1.m3u8
四川影视文艺,http://scgctvshow.sctv.com/hdlive/sctv5/index.m3u8
四川文化旅游,http://scgctvshow.sctv.com/hdlive/sctv2/index.m3u8
四川文化旅游,http://scgctvshow.sctv.com/hdlive/sctv2/1.m3u8
四川新闻,http://scgctvshow.sctv.com/hdlive/sctv4/index.m3u8
四川新闻频道,http://scgctvshow.sctv.com/hdlive/sctv4/1.m3u8
四川星空购物,http://scgctvshow.sctv.com/hdlive/sctv6/index.m3u8
四川科教,http://3017ugjo.live2.danghongyun.com/live/hls/4dc0e927160647a2bd585cb5495810be/3c2126c925204cddb74a3534b97e2765-1.m3u8
四川经济,http://scgctvshow.sctv.com/hdlive/sctv3/index.m3u8
四川乡村 http://m3u8.sctv.com/tvlive/SCTV9/index.m3u8
四川卫视 http://m3u8.sctv.com/tvlive/SCTV0/index.m3u8
四川科教 http://m3u8.sctv.com/tvlive/SCTV8/index.m3u8
四川康巴卫视 Kangba
经济 2
http://www.sctv.com/

 */
func init()  {
	addSctv("SCWS", "四川卫视", "SCTV0", "卫视","720")
	addSctv("SCKJ", "四川科教", "SCTV8", "科教","720")
	addSctv("SCWHLY", "四川文化旅游", "SCTV3", "文旅","720")
	addSctv("SCJJ", "四川经济", "SCTV2", "财经","720")
	addSctv("SCYSWH", "四川新闻", "SCTV4", "新闻","720")
	addSctv("SCYSWH", "四川影视文艺", "SCTV5", "影视","720")
	addSctv("XKGW", "四川星空购物", "SCTV6", "购物","720")
	addSctv("SCFNET", "四川妇女儿童", "SCTV7", "少儿","720")
	addSctv("SCKBWS", "四川康巴卫视", "Kangba", "卫视","720")


	addCdtv("CDXWZH", "成都新闻综合", "1", "综合","普清")
	addCdtv("CDJJZX", "成都经济咨询", "2", "财经","普清")
	addCdtv("CDDSSH", "成都都市生活", "3", "综合","普清")
	addCdtv("CDYSWY", "成都影视文艺", "4", "影视","普清")
	addCdtv("CDGG", "成都公共", "5", "综合","普清")
	addCdtv("CDSE", "成都少儿", "6", "少儿","普清")
}

type IPTvSource struct {
	Name string
	Type string
	Host string
	Path string
}


func addSctv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Group: "四川",
		Category: category,
		Quality: quality,
		Name:    name,
		DirectReturn: false,
		Referer: "http://www.sctv.com/",
		ReRegxp: reRegx2,
		Prefix: "http://m3u8.sctv.com/tvlive/%s/",
		UrlFmt:  "http://m3u8.sctv.com/tvlive/%s/index.m3u8",
	}
}


func addCdtv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Group: "四川",
		Name:    name,
		Referer: "https://www.cditv.cn/",
		UrlFmt:  "https://www.cditv.cn/api.php?op=live&type=playTv&fluency=sd&videotype=m3u8&catid=192&id=%s",
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(string(body), 302).ServeHTTP(w, r)
		},
	}
}
