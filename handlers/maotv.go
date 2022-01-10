package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"wmenjoy.com/iptv/utils"
)

//https://pastebin.com/raw/KGRduBqa
//https://pastebin.com/raw/KGRduBqa
//	parserList = append(parserList, videoParserInfo{
//		Name: "4K影院1",
//		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=1095368&key=cdgjpsuDFLNRTUVY03&url=%s",
//	})
//	parserList = append(parserList, videoParserInfo{
//		Name: "绿箭影视2",
//		Url:  "https://json.pangujiexi.com:12345/json.php?url=%s",
//	})

type MaoTVLiveChannel struct {
	Name string `json:"name"`
	Urls []string `json:"urls"`

}
type MaoTVLive struct {
	Group string `json:"group"`
	Channels []MaoTVLiveChannel `json:"channels"`
}
type MaoTVSite struct {
	Key string `json:"key"`
	Name string `json:"name"`
	Type int `json:"type"`
	Api string `json:"api"`
	Searchable int `json:"searchable"`
	QuickSearch string `json:"quickSearch"`
	Filterable string `json:"filterable"`
}

type MaoTVParser struct {
	Name string `json:"name"`
	Type int `json:"type"`
	Url string `json:"url"`
}

type MaoTVOption struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Category int `json:"category"`

}
type MaoTVIJK struct {
	Group string `json:"group"`
	Options []MaoTVOption `json:"options"`
}
type MaoTVInfo struct {
	Sites []MaoTVSite `json:"sites"`
	Lives []MaoTVLive `json:"lives"`
	Parses []MaoTVParser `json:"parses"`
	Flags []string `json:"flags"`
	Ijk	 []MaoTVIJK `json:"ijk"`
	Ads  []string `json:"ads"`
	wallpaper string `json:"wallpaper"`
	spider string `json:"spider"'`
}

/**
parserList = append(parserList, videoParserInfo{
	Name: "追剧吧2",
	Url: "https://svip.spchat.top/api/?key=SAl7tLs3Zzu5alSNtz&url=%s",
})
 */
type videoParserResult struct {
	Code        string `json:"code"`
	Success     string `json:"success"`
	Type        string `json:"type"`
	Player      string `json:"player"`
	Msg         string `json:"msg"`
	Url         string `json:"url"`
	DaluIDC     string `json:"DaluIDC"`
	HongKongIDC string `json:"HongKongIDC"`
	Txt         string `json:"txt"`
	Txt2        string `json:"txt2"`
	From        string `json:"From"`
	FromUrl     string `json:"From_Url"`
}

type videoParserInfo struct {
	Name string
	Url  string
	Replace bool
}

var parserList = make([]videoParserInfo, 0)

func init() {

	parserList = append(parserList, videoParserInfo{
		Name: "黄河影视",
		Url:  "http://jx.ledu8.cn/api/?key=P8QSgO61p1MpHV2ALn&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "瑞丰资源",
		Url:  "https://api.m3u8.tv:5678/home/api?type=ys&uid=1285201&key=bcikqtwxyADEGKUX36&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "筋斗云",
		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=1494542&key=ijmqvwxABEHILMOT48&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "久九影视",
		Url: "https://json.hfyrw.com/mao.go?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "蓝光视频",
		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=1931000&key=gktuvyzABEORSYZ135&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "益达影院",
		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=123503&key=fgkqryzDEFLNQSTW69&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "(无名)",
		Url: "https://api.zakkpa.com:8888/analysis/json/?uid=91&my=gjksuvCHIJLRS01268&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "DC影视",
		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=7665652&key=dglmnwEFILMOPRW056&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "抹茶猪",
		Url: "https://api.m3u8.tv:5678/home/api?type=ys&uid=1525223&key=fhikpsvBCDFHJOSUZ8&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "麻瓜视频",
		Url: "https://api.parwix.com:4433/analysis/json/?uid=1735&my=cejkmnuvyBEFINR056&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "Yoyo1",
		Url: "https://json.legendwhb.cn/json.php/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "爱影视",
		Url:  "http://14.17.115.200:520/json.php?id=6e5LaYyU5JLs9aRawyGKwPkH7ZFr701z&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "追剧吧1",
		Replace: true,
		Url:  "http://newjiexi.gotto.top/yun_apib.php/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "优视影视",
		Replace: true,
		Url: "http://ccc.ysys.asia/jx/?id=2&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "黄河影视",
		Url:  "http://jx.ledu8.cn/api/?key=P8QSgO61p1MpHV2ALn&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "(无名)1",
		Url:  "http://vipjh.chunbaotaiji.com/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "麻瓜视频1",
		Url:  "https://jhjx.ptygx.com/xttyjx.php/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "Vip影院",
		Url:  "https://api.zakkpa.com:8888/analysis/json/?uid=39&my=ehklrtxzAFKLMNUXY2&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "大头影视",
		Url:  "http://fast.rongxingvr.cn:8866/api/?key=M3tZzS2q0oGrQ7aWlr&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "绿箭影视",
		Url:  "http://jf.jisutuku.top/api/?key=RHjXcjUTkyZnfWx9u4&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "Yoyo",
		Replace: true,
		Url:  "http://47.100.138.210:91/home/api?type=ys&uid=1947441&key=bcfgjmuwCEORSVX237&url=%s",
	})

	parserList = append(parserList, videoParserInfo{
		Name: "蜜果TV",
		Url:  "https://fast.rongxingvr.cn:8866/api/?key=rmxOw7BpINuGIyWQng&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "筋斗云1",
		Url:  "https://fast.rongxingvr.cn:8866/api/?key=nrR7koAyq9ajKId4nC&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "麻花视频",
		Url:  "http://jf.tcspvip.com:246/home/api?type=ys&uid=65404&key=bjoprtvyABGIMPXZ27&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "TV酷",
		Url:  "http://a.dxzj88.com/jhjson/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "饭后电影",
		Url:  "https://fast.rongxingvr.cn:8866/api/?key=J4mUIu3DrRtIOojDox&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "影阅阁2",
		Url:  "https://api.m3u8.tv:5678/home/api?type=ys&uid=594615&key=bcehpqtxCEGKMT0248&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "豆渣影视",
		Url:  "https://fast.rongxingvr.cn:8866/api/?key=jtDZ22biNujOBLlgoe&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "看剧吧",
		Url:  "http://47.100.138.210:91/home/api?type=ys&uid=243653&key=kqswxyABGHKLQSV127&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "玺娜影视",
		Url: "https://balabala.yatongle.com/api/?key=RuHZpg9zxigiLZRIyl&url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "遗忘影视",
		Url: "http://gq.bywdtk.cn/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "星空影视",
		Url: "http://api.vip123kan.vip/?url=%s",
	})
	parserList = append(parserList, videoParserInfo{
		Name: "爱追剧",
		Url: "https://jhjx.ptygx.com/xttyjx.php/?url=%s",
	})

	go func() {
		temp := make([]videoParserInfo, 0)
		for _, v := range parserList {
			start := time.Now()
			video := getRealVedioUrl("https://v.qq.com/x/cover/mzc00200qqsk3cv.html", v.Url, v.Replace)
			if video == "" {
				continue
			}
			tc:=time.Since(start)
			if  tc < time.Duration(500000000) {
				temp = append(temp, v)
			}
			fmt.Printf("%s time cost = %v\n", v.Url, tc)

		}
		parserList = temp
	}()


}
func timeCost() func() {
	start := time.Now()
	return func() {
		tc:=time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}

func getRealVedioUrl(src string, parserUrl string, transfer bool)string{
	realUrl := fmt.Sprintf(parserUrl, src)

	req, err := http.NewRequest("GET", realUrl, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}


	videoItem := utils.MatchOneOf(string(body), "\"url\"[ ]*:[ ]*\"([^\"]+)")

	if videoItem == nil || len(videoItem) <= 1 {
		print(parserUrl)
		return ""
	}

	videoUrl := strings.ReplaceAll(videoItem[1], "\\", "")

	return videoUrl
}
func MaoTvHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	url := r.Form.Get("url")
	if url == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	for i := 0;i < 5; i++{
		index := rand.Intn(len(parserList) -1 )
		parser := parserList[index]
		videoUrl := getRealVedioUrl(url, parser.Url, parser.Replace)
		if videoUrl != "" {
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(videoUrl, 302).ServeHTTP(w, r)
			break
		}
	}

}
