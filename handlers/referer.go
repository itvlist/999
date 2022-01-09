package handlers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wmenjoy.com/iptv/utils"
)

var host = utils.GetIP()

type RerferInfo struct {
	Id string
	Referer string
	urlFmt string
	reRegxp *regexp.Regexp
	prefix string
	urlBuildFunc func(refererInfo RerferInfo) string
	beforeFunc func(refererInfo RerferInfo, url string, header http.Header)
	afterFunc func(refererInfo RerferInfo, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request)
}
var reRegx2 = regexp.MustCompilePOSIX(`([^#]+\.ts)`)
/***
 *
 <?php
$id = isset($_GET['id']) ? $_GET['id'] : 'cetv1';
$n = array(
    'cetv1' => 695,
    'cetv2' => 696,
    'cetv3' => 697,
    'cetv4' => 698,
);
$ch = curl_init('http://app.cetv.cn/video/player/stream?site_id=10001&stream_id=' . $n[$id]);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$data = curl_exec($ch);
curl_close($ch);
$playUrl = json_decode($data)->stream;
$pathArr = explode("/", $playUrl);
$endPath = $pathArr[sizeof($pathArr) - 1];
$endPath = str_replace(".m3u8", "", $endPath);
$c = curl_init($playUrl);
curl_setopt($c, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($c, CURLOPT_HTTPHEADER, array("Referer: http://app.cetv.cn"));
$result = curl_exec($c);
curl_close($c);
$ts = str_replace($endPath, "http://txycsbl.centv.cn/zb/" . $endPath, $result);
print_r($ts);
?>
 */
var referinfos = make(map[string]*RerferInfo,0)


func addCetv(name string, id string){
	referinfos[name] = &RerferInfo{
		Id: id,
		Referer:"http://app.cetv.cn",
		urlFmt: "http://txycsbl.centv.cn/zb/0104%s.m3u8",
		prefix: "http://txycsbl.centv.cn/zb/",
		reRegxp:  regexp.MustCompilePOSIX(`([^#]+\.ts)`),
	}
}

type nodeParamResult struct {
	Node string `json:"node"`
}

type gdtvQueryResult struct {
	AvatarUrl string `json:"avatarUrl"`
	Category int `json:"category"`
	CoverUrl string `json:"coverUrl"`
	Keyword string `json:"keyword"`
	Name string `json:"name"`
	Pk int `json:"pk"`
	PlayUrl string `json:"playUrl"`
	Slogan string `json:"slogan"`
	TimeOffset int `json:"timeOffset"`

}
type gdtvPlayUrlResult struct {
	Hd string `json:"hd"`
}


func optionForUrl(url string){
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("access-control-request-method", "GET")
	req.Header.Set("authority", "gdtv-api.gdtv.cn")
	req.Header.Set("access-control-request-headers", "content-type,x-itouchtv-ca-key,x-itouchtv-ca-signature,x-itouchtv-ca-timestamp,x-itouchtv-client,x-itouchtv-device-id")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("origin", `https://www.gdtv.cn`)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

}

func addGdtv(name string, id string){
	referinfos[name] = &RerferInfo{
		Id: id,
		Referer:"http://app.cetv.cn",
		urlFmt: "https://gdtv-api.gdtv.cn/api/tv/v2/tvChannel/%s?tvChannelPk=%s&node=%s",
		prefix: "http://txycsbl.centv.cn/zb/",
		reRegxp:  regexp.MustCompilePOSIX(`([^#]+\.ts)`),
		urlBuildFunc: func(refererInfo RerferInfo) string {
			req, err := http.NewRequest("GET", "https://tcdn-api.itouchtv.cn/getParam", nil)
			node := ""
			if err != nil {
				node ="2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				node ="2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				node ="2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			paramResult := nodeParamResult{}
			err = json.Unmarshal(body, &paramResult)
			if err != nil || paramResult.Node == "" {
				node ="2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			} else {
				node = paramResult.Node
			}

			url := fmt.Sprintf(refererInfo.urlFmt, id, id, base64.StdEncoding.EncodeToString([]byte(node)))
			optionForUrl(url)
			return url
		},
		beforeFunc: func(refererInfo RerferInfo, url string, header http.Header) {
			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
			header.Set("authority", "gdtv-api.gdtv.cn")
			header.Set("x-itouchtv-ca-timestamp", timestamp)
			header.Set("x-itouchtv-ca-key", "89541443007807288657755311869534")
			header.Set("x-itouchtv-client", "WEB_PC")
			header.Set("x-itouchtv-device-id", "WEB_b9cd0e90-709e-11ec-8f23-ebd518a0a397")
			header.Set("content-type", "application/json")
			header.Set("accept", "application/json, text/plain, */*")
			header.Set("origin", "https://www.gdtv.cn")
			header.Set("sec-ch-ua-mobile", "?0")
			header.Set("sec-ch-ua-platform", "macOS")
			header.Set("sec-fetch-site", "same-site")
			header.Set("sec-fetch-mode", "cors")
			header.Set("accept-language", "zh-CN,zh;q=0.9")
			secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
			message := []byte(fmt.Sprintf("GET\n%s\n%s\n",url, timestamp))
			hash := hmac.New(sha256.New, secret)
			hash.Write(message)
			signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			header.Set("x-itouchtv-ca-signature", signature)
			header.Set("Referer", "https://www.gdtv.cn/")

		},
		afterFunc: func(refererInfo RerferInfo, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}
			gdtvQueryResult := gdtvQueryResult{}
			err = json.Unmarshal(body, &gdtvQueryResult)
			gdtvPlayUrlResult := gdtvPlayUrlResult{}
			err = json.Unmarshal([]byte(gdtvQueryResult.PlayUrl), &gdtvPlayUrlResult)

			url := fmt.Sprintf("http://%s:8880/transfer?url=%s&referer=%s",host,
				base64.StdEncoding.EncodeToString([]byte(gdtvPlayUrlResult.Hd)),base64.StdEncoding.EncodeToString([]byte("https://www.gdtv.cn/")))
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(url,302).ServeHTTP(w, r)


			/**

			*/
		},
	}
}

func addSztv(name string, id string){
	referinfos[name] = &RerferInfo{
		Id:      id,
		urlFmt:  "https://sztv-live.cutv.com/%s/%s/%s.m3u8",
		reRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
		urlBuildFunc: func(refererInfo RerferInfo) string {
			secret := "cutvLiveStream|Dream2017"
			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:10]
			h := md5.New()
			h.Write([]byte(timestamp + id + secret))
			a := h.Sum(nil)
			c := hex.EncodeToString(a)

			numUrl := fmt.Sprintf(
				"https://cls2.cutv.com/getCutvHlsLiveKey?t=%s&id=%s&token=%s&at=1", timestamp, id, c);
			req, err := http.NewRequest("GET", numUrl, nil)
			if err != nil {
				return ""
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}

			strBody := strings.Trim(string(body), "\"")
			fileName := ""
			index := len(strBody)
			if (index > 0) {
				index -= index / 2
				strBody = strBody[index:] + strBody[0:index]
				fileNameb, _ := base64.StdEncoding.DecodeString(reverseString(strBody))
				fileName = string(fileNameb)
			}

			url := fmt.Sprintf(refererInfo.urlFmt, id, "500", fileName)
			return url
		},

	}
}

func reverseString(str string) string{
	runes := []rune(str)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)

}

func m(){
	req, err := http.NewRequest("GET"," gdtvPlayUrlResult.Hd", nil)

	if err != nil {
		//http.Error(w, err.Error(), 503)
		return
	}

	req.Header.Set("authority", "tcdn.itouchtv.cn")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Referer", "https://www.gdtv.cn/")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("accept", `*/*`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("origin", `https://www.gdtv.cn`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	//	http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//http.Error(w, err.Error(), 503)
		return
	}
	//w.Write(body)
}

func init(){
	addCetv("CETV-1","cetv1")
	addCetv("CETV-2","cetv2")
	addCetv("CETV-3","cetv3")
	addCetv("CETV-4","cetv4")
	addGdtv("GDWS","43")
	addGdtv("GDZJ","44")
	addGdtv("GDXW","45")
	addGdtv("GDTY","47")
	addGdtv("NFWS","51")
	addGdtv("GDJJKJ","49")
	addGdtv("GDYS","53")
	addGdtv("GDZY","16")
	addGdtv("GDGJ","46")
	addGdtv("GDSE","54")
	addGdtv("GDJJKT","66")
	addGdtv("GDNFGW","42")
	addGdtv("GDLNXQ","15")
	addGdtv("GDFC","67")
	addGdtv("GDXDJY","13")
	addGdtv("GDYD","74")
	addGdtv("GRTNWHPD","75")

	addSztv("SZWS","AxeFRth")
	addSztv("SZYL","1q4iPng")
	addSztv("SZSE","1SIQj6s")
	addSztv("SZGG","2q76Sw2")
	addSztv("SZDSJ","4azbkoY")
	addSztv("SZDB","9zoW71b")
	addSztv("SZYHGW","BJ5u5k2")
	addSztv("SZDS","ZwxzUXr")
	addSztv("SZGJ","sztvgjpd")
	addSztv("SZGJ","sztvgjpd")
	addSztv("SZTYJK","sztvtyjk")
	addSztv("SZLG","uGzbXhS")
	addSztv("SZYDSX","wDF6KJ3")
	addSztv("SZDVSH","xO1xQFv")

}

func TransferHandler(w http.ResponseWriter, r *http.Request){
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url := r.Form.Get("url")
	referer := r.Form.Get("referer")
	urlb ,err := base64.StdEncoding.DecodeString(url);
	refererb, err:=  base64.StdEncoding.DecodeString(referer);

	req, err := http.NewRequest("GET",string(urlb), nil)

	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Referer", string(refererb))
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("accept", `*/*`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	w.Write(body)
}


func RefererHandler(w http.ResponseWriter, r *http.Request)  {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url := ""

	id := r.Form.Get("id")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	referinfo := referinfos[strings.ToUpper(id)]

	if referinfo == nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	if referinfo.urlBuildFunc == nil {
		url = fmt.Sprintf(referinfo.urlFmt, referinfo.Id)
	} else {
		url = referinfo.urlBuildFunc(*referinfo)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if referinfo.beforeFunc != nil {
		referinfo.beforeFunc(*referinfo, url, req.Header)
	} else {
		if referinfo.Referer != ""{
			req.Header.Set("Referer", referinfo.Referer)
		}
		req.Header.Set("accept", `*/*`)
	}



	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	if referinfo.afterFunc == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
		if referinfo.prefix != "" {
			newbody := reRegx.ReplaceAll(body, []byte(referinfo.prefix+"$0"))

			logrus.Info(string(newbody))
			w.Write(newbody)
		} else {
			w.Write(body)
		}
	} else {
		referinfo.afterFunc(*referinfo, resp.Body, w, r)
	}

}
