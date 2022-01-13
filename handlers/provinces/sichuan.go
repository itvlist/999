package provinces

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
四川·康巴卫视,http://scgctvshow.sctv.com:80/hdlive/kangba/1.m3u8
四川·康巴卫视,http://scgctvshow.sctv.com/hdlive/kangba/1.m3u8
四川公共乡村,http://scgctvshow.sctv.com/hdlive/sctv9/index.m3u8
四川公共乡村,http://scgctvshow.sctv.com/hdlive/sctv9/1.m3u8
四川乡村 http://m3u8.sctv.com/tvlive/SCTV9/index.m3u8
四川卫视 http://m3u8.sctv.com/tvlive/SCTV0/index.m3u8
四川科教 http://m3u8.sctv.com/tvlive/SCTV8/index.m3u8
四川康巴卫视 Kangba
经济 2
http://www.sctv.com/

*/
/**
成都新闻,http://ye23.win/iptv/cdtvhls.php?id=cdxwzh
成都新闻,http://ye23.win/iptv/cdtvflv.php?id=cdxwzh
'cdxwzh' => 1,//成都新闻综合
'cdjjzx' => 2,//成都经济资讯
'cddssh' => 3,//成都都市生活
'cdyswy' => 45,//成都影视文艺
'cdgg' => 5,//成都公共
'cdse' => 6,//成都少儿
'cdss' => 9,//成都时尚
'cdqc' => 10,//成都汽车
'cdzxxgx' => 11,//成都资讯新干线
'cdmsly' => 12,//成都美食旅游
'cdrcxf' => 15,//成都蓉城先锋
'cdmrgw' => 18,//成都每日购物

'cdgjyd' => 1152,//成都公交移动*
'jygw' => 560,//家有购物*

'jntv' => 556,//金牛电视台
'slzh' => 557,//双流综合
'wjtv' => 559,//温江电视台
'gxtv' => 722,//高新电视台
'xjtv' => 760,//新津电视台
'dyxwzh' => 790,//大邑综合
'pztv' => 796,//彭州电视台
'pjtv' => 828,//蒲江电视台
'jjtv' => 1541,//锦江电视台
'jttv' => 840,//金堂电视台
'pdxwzh' => 845,//郫都新闻综合
'lqzh' => 882,//龙泉综合
'qytv' => 910,//青羊电视台
'qbjtv' => 966,//青白江电视台
'cztv' => 1257,//崇州电视1套
'djytv' => 1314,//都江堰电视台
'chyx' => 1319,//成华有线
'qltv' => 1427,//邛崃电视台
'whtv' => 1766,//武侯电视台
'xdzh' => 1712,//新都电视台
'jyxwzh' => 1698,//简阳新闻综合

*/
func init() {
	addSctv("SCWS", "四川卫视", "SCTV0", "卫视", "720")
	addSctv("SCKJ", "四川科教", "SCTV8", "科教", "720")
	addSctv("SCWHLY", "四川文化旅游", "SCTV3", "文旅", "720")
	addSctv("SCJJ", "四川经济", "SCTV2", "财经", "720")
	addSctv("SCYSWH", "四川新闻", "SCTV4", "新闻", "720")
	addSctv("SCYSWH", "四川影视文艺", "SCTV5", "影视", "720")
	addSctv("XKGW", "四川星空购物", "SCTV6", "购物", "720")
	addSctv("SCFNET", "四川妇女儿童", "SCTV7", "少儿", "720")
	addSctv("SCKBWS", "四川康巴卫视", "Kangba", "卫视", "720")
	addSctv("SCXC", "四川公共乡村", "SCTV9", "农民", "720")
	addCdtv("CDXWZH", "成都新闻综合", "1", "综合", "普清")
	addCdtv("CDJJZX", "成都经济咨询", "2", "财经", "普清")
	addCdtv("CDDSSH", "成都都市生活", "3", "综合", "普清")
	addCdtv("CDYSWY", "成都影视文艺", "4", "影视", "普清")
	addCdtv("CDGG", "成都公共", "5", "综合", "普清")
	addCdtv("CDSE", "成都少儿", "6", "少儿", "普清")
}

type IPTvSource struct {
	Name string
	Type string
	Host string
	Path string
}

func addSctv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:           id,
		Key:          key,
		Jump:         true,
		Group:        "四川",
		Category:     category,
		Quality:      quality,
		Name:         name,
		DirectReturn: false,
		Referer:      "http://www.sctv.com/",
		ReRegxp:      reRegx2,
		Prefix:       "http://m3u8.sctv.com/tvlive/%s/",
		UrlFmt:       "http://m3u8.sctv.com/tvlive/%s/index.m3u8",
	}
}

func addCdtv(key string, name string, id string, category string, quality string) {
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:      id,
		Key:     key,
		Jump:    true,
		Group:   "四川",
		Name:    name,
		Referer: "https://www.cditv.cn/",
		UrlFmt:  "https://www.cditv.cn/api.php?op=live&type=playTv&fluency=sd&videotype=m3u8&catid=192&id=%s",
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl;charset=UTF-8")
			http.RedirectHandler(string(body), 302).ServeHTTP(w, r)
		},
	}
}
