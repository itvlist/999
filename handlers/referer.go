package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/sirupsen/logrus"
	"wmenjoy.com/iptv/utils"
)

var host = utils.GetIP()

type RerferInfo struct {
	Id           string
	Key          string
	Name         string
	Referer      string
	urlFmt       string
	Jump         bool
	reRegxp      *regexp.Regexp
	prefix       string
	urlBuildFunc func(refererInfo RerferInfo) string
	beforeFunc   func(refererInfo RerferInfo, url string, header http.Header)
	afterFunc    func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request)
}

var reRegx2 = regexp.MustCompilePOSIX(`([^#]+\.ts)`)
var referinfos = make(map[string]*RerferInfo, 0)

func addCetv(name string, id string) {
	referinfos[name] = &RerferInfo{
		Id:      id,
		Referer: "http://app.cetv.cn",
		urlFmt:  "http://txycsbl.centv.cn/zb/0104%s.m3u8",
		prefix:  "http://txycsbl.centv.cn/zb/",
		reRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
	}
}

type nodeParamResult struct {
	Node string `json:"node"`
}

type gdtvQueryResult struct {
	AvatarUrl  string `json:"avatarUrl"`
	Category   int    `json:"category"`
	CoverUrl   string `json:"coverUrl"`
	Keyword    string `json:"keyword"`
	Name       string `json:"name"`
	Pk         int    `json:"pk"`
	PlayUrl    string `json:"playUrl"`
	Slogan     string `json:"slogan"`
	TimeOffset int    `json:"timeOffset"`
}
type gdtvPlayUrlResult struct {
	Hd string `json:"hd"`
}


func optionForUrl(url string) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
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

//iv'1563432177954301' pkpadding7 004BC54E474B0F72902CE2E29B91C5E0A3E9E3BA11435370FAEAC96B25C9805D
func optionForGetParam(url string) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("access-control-request-method", "GET")
	req.Header.Set("accept", "*/*")
	req.Header.Set("authority", "tcdn-api.itouchtv.cn")
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

func addGdtv(name string, id string) {
	referinfos[name] = &RerferInfo{
		Id:      id,
		Referer: "http://app.cetv.cn",
		urlFmt:  "https://gdtv-api.gdtv.cn/api/tv/v2/tvChannel/%s?tvChannelPk=%s&node=%s",
		prefix:  "http://txycsbl.centv.cn/zb/",
		reRegxp: regexp.MustCompilePOSIX(`([^#]+\.ts)`),
		urlBuildFunc: func(refererInfo RerferInfo) string {
			optionForGetParam("https://tcdn-api.itouchtv.cn/getParam")
			req, err := http.NewRequest("GET", "https://tcdn-api.itouchtv.cn/getParam", nil)
			node := ""
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}

			timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
			req.Header.Set("authority", "tcdn-api.itouchtv.cn")
			req.Header.Set("x-itouchtv-ca-timestamp", timestamp)
			req.Header.Set("x-itouchtv-ca-key", "89541443007807288657755311869534")
			req.Header.Set("x-itouchtv-client", "WEB_PC")
			req.Header.Set("x-itouchtv-device-id", "WEB_d547fdf0-633e-11ec-83d3-fb13b5511434")
			req.Header.Set("content-type", "application/json")
			req.Header.Set("accept", "application/json, text/plain, */*")
			req.Header.Set("origin", "https://www.gdtv.cn")
			req.Header.Set("sec-ch-ua-mobile", "?0")
			req.Header.Set("sec-ch-ua-platform", "macOS")
			req.Header.Set("sec-fetch-site", "same-site")
			req.Header.Set("sec-fetch-mode", "cors")
			req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
			secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
			message := []byte(fmt.Sprintf("GET\n%s\n%s\n", "https://tcdn-api.itouchtv.cn/getParam", timestamp))
			hash := hmac.New(sha256.New, secret)
			hash.Write(message)
			signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			req.Header.Set("x-itouchtv-ca-signature", signature)
			req.Header.Set("Referer", "https://www.gdtv.cn/")
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
			}
			paramResult := nodeParamResult{}
			err = json.Unmarshal(body, &paramResult)
			if err != nil || paramResult.Node == "" {
				node = "2028330049-a01d3e9adb2e3ef55a6b897cf943e947"
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
			header.Set("x-itouchtv-device-id", "WEB_d547fdf0-633e-11ec-83d3-fb13b5511434")
			header.Set("content-type", "application/json")
			header.Set("accept", "application/json, text/plain, */*")
			header.Set("origin", "https://www.gdtv.cn")
			header.Set("sec-ch-ua-mobile", "?0")
			header.Set("sec-ch-ua-platform", "macOS")
			header.Set("sec-fetch-site", "same-site")
			header.Set("sec-fetch-mode", "cors")
			header.Set("accept-language", "zh-CN,zh;q=0.9")
			secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
			message := []byte(fmt.Sprintf("GET\n%s\n%s\n", url, timestamp))
			hash := hmac.New(sha256.New, secret)
			hash.Write(message)
			signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			header.Set("x-itouchtv-ca-signature", signature)
			header.Set("Referer", "https://www.gdtv.cn/")

		},
		afterFunc: func(refererInfo RerferInfo, srcUrl string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
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

			values := url.Values{}
			values.Set("url", gdtvPlayUrlResult.Hd)
			values.Set("referer", "https://www.gdtv.cn/")

			urlstr := fmt.Sprintf("http://%s:8880/transfer?%s", host, values.Encode())
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(urlstr, 302).ServeHTTP(w, r)
		},
	}
}

func addSztv(name string, id string) {
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
				"https://cls2.cutv.com/getCutvHlsLiveKey?t=%s&id=%s&token=%s&at=1", timestamp, id, c)
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
			if index > 0 {
				index -= index / 2
				strBody = strBody[index:] + strBody[0:index]
				fileNameb, _ := base64.StdEncoding.DecodeString(reverseString(strBody))
				fileName = string(fileNameb)
			}

			url := fmt.Sprintf(refererInfo.urlFmt, id, "500", fileName)
			return url
		},
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()

			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}

			index := strings.LastIndex(url, "/")
			newbody := reRegx.ReplaceAll(body, []byte(url[0:index]+"/$0"))

			logrus.Info(string(newbody))
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			w.Write(newbody)
		},
	}
}

func addCqtv(key string, name string, id string) {
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Name:   name,
		urlFmt: "https://web.cbg.cn/live/getLiveUrl?url=%s",
		urlBuildFunc: func(refererInfo RerferInfo) string {

			url := fmt.Sprintf("https://rmtapi.cbg.cn/list/%s/1.html?pagesize=20", refererInfo.Id)
			req, err := http.NewRequest("GET", url, nil)
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
			urlResult := utils.MatchOneOf(string(body), "\"ios_HDlive_url\":\"([^\"]+)")[1]
			return fmt.Sprintf(refererInfo.urlFmt, urlResult)
		},
		afterFunc: func(refererInfo RerferInfo, srcUrl string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}
			urlResult := utils.MatchOneOf(string(body), "\"url\"[ ]*:[ ]*\"([^\"]+)")[1]
			proxyUrl := fmt.Sprintf("http://%s:8880/ats?file=", host)
			index := strings.LastIndex(urlResult, "/")
			values := url.Values{}
			values.Set("url", urlResult)
			values.Set("srcUrl", urlResult[0:index])
			values.Set("proxy", proxyUrl)
			values.Set("key", getKey())
			values.Set("referer", "https://www.cbg.cn/")
			realUrl := fmt.Sprintf("http://%s:8880/transfer?%s", host, values.Encode())
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(realUrl, 302).ServeHTTP(w, r)
		},
	}
}

func getKey() string {
	url := "https://sjlivecdnx.cbg.cn/1ive/stream_2.php"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("referer", "https://www.cbg.cn/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(body)
}

type xjParamInfo struct {
	Paramslist struct {
		Language           string `json:"language"`
		SkinType           string `json:"skinType"`
		PlayerId           string `json:"playerId"`
		StreamType         string `json:"streamType"`
		Logging            bool   `json:"logging"`
		LogLevel           string `json:"logLevel"`
		LogFilter          string `json:"logFilter"`
		Volume             int    `json:"volume"`
		Loop               bool   `json:"loop"`
		Smoothing          bool   `json:"smoothing"`
		AutoPlay           bool   `json:"autoPlay"`
		AutoLoad           bool   `json:"autoLoad"`
		Skin               string `json:"skin"`
		BufferTime         int    `json:"bufferTime"`
		Configable         bool   `json:"configable"`
		Host               string `json:"host"`
		Version            string `json:"version"`
		PlayerBackground   string `json:"playerBackground"`
		Plugin             bool   `json:"plugin"`
		NonDisplay         string `json:"nonDisplay"`
		SeekParam          string `json:"seekParam"`
		TimeServer         string `json:"timeServer"`
		Programchanggehost string `json:"programchanggehost"`
		AudioOnly          bool   `json:"audioOnly"`
		Isshowcontrol      bool   `json:"isshowcontrol"`
		IsUrlStatic        bool   `json:"isUrlStatic"`
		EncryptUrl         bool   `json:"EncryptUrl"`
		EncryptionSwf      string `json:"EncryptionSwf"`
		Loadinglogo        string `json:"loadinglogo"`
		Isfullscreen       bool   `json:"isfullscreen"`
		Ispause            bool   `json:"ispause"`
	} `json:"paramslist"`
	Pluginslist []struct {
		Source   string `json:"source"`
		Callback string `json:"callback"`
	} `json:"pluginslist"`
	ParamsConfig struct {
		CdnConfig []struct {
			Code             string `json:"code"`
			PublishHost      string `json:"publishHost"`
			H5PublishHost    string `json:"H5PublishHost"`
			SeekField        string `json:"seekField"`
			Unit             string `json:"unit"`
			OpenChain        string `json:"openChain"`
			InvalidTime      string `json:"invalidTime"`
			EncryptMode      string `json:"encryptMode"`
			OpenPcdn         string `json:"openPcdn"`
			GetAuthUrl       string `json:"getAuthUrl"`
			PlaybackLiveHost string `json:"PlaybackLiveHost"`
		} `json:"cdnConfig"`
		CdnConfigEncrypt   string `json:"cdnConfigEncrypt"`
		IsCDNConfigEncrypt bool   `json:"isCDNConfigEncrypt"`
	} `json:"paramsConfig"`
}
type xjEncryptInfo struct {
	Code             string `json:"code"`
	PublishHost      string `json:"publishHost"`
	H5PublishHost    string `json:"H5PublishHost"`
	SeekField        string `json:"seekField"`
	Unit             string `json:"unit"`
	OpenChain        string `json:"openChain"`
	InvalidTime      string `json:"invalidTime"`
	EncryptKey       string `json:"encryptKey"`
	LiveEncryptKey   string `json:"liveEncryptKey"`
	EncryptMode      string `json:"encryptMode"`
	OpenPcdn         string `json:"openPcdn"`
	GetAuthUrl       string `json:"getAuthUrl"`
	PlaybackLiveHost string `json:"PlaybackLiveHost"`
}

func addXjtv(key string, name string, id string) {
	PKCS7UnPadding := func(origData []byte) []byte {
		length := len(origData)
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)]
	}
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		prefix: "http://livehyw5.chinamcache.com/hyw/",
		urlFmt: "http://livehyw5.chinamcache.com/hyw/%s.m3u8?txSecret=%s&txTime=%s",
		urlBuildFunc: func(refererInfo RerferInfo) string {
			return "http://mediaxjtvs.chinamcache.com/hyw/media/playerJson/liveChannel/7d40edeb62fe4f8a9d9a08bc653dcab6_PlayerParamProfile.json"
		},
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			defer bodyReader.Close()

			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, err.Error(), 503)
				return
			}

			paramInfo := xjParamInfo{}
			json.Unmarshal(body, &paramInfo)
			block, err := aes.NewCipher([]byte("roZ68Okc5MUTMraM"))
			result, _ := hex.DecodeString(paramInfo.ParamsConfig.CdnConfigEncrypt)
			//blockSize := block.BlockSize()
			//tmp := ZeroPadding([]byte(data), blockSize)
			mode := cipher.NewCBCDecrypter(block, []byte("7384627385960726"))
			text := make([]byte, len(result))
			mode.CryptBlocks(text, result)

			encInfoList := make([]xjEncryptInfo, 0)
			value := PKCS7UnPadding(text)
			json.Unmarshal(value, &encInfoList)
			invalidTime, _ := strconv.ParseInt(encInfoList[0].InvalidTime, 10, 64)
			timestamp := strconv.FormatInt(time.Now().Unix()+invalidTime, 16)
			h := md5.New()
			h.Write([]byte(encInfoList[0].EncryptKey + refererInfo.Id + timestamp))
			a := h.Sum(nil)
			realUrl := fmt.Sprintf(refererInfo.urlFmt, refererInfo.Id, hex.EncodeToString(a), timestamp)
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(realUrl, 302).ServeHTTP(w, r)
		},
	}
}


type JxtvAuthResult struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	T      string `json:"t"`
	Count  int    `json:"count"`
	Token  string `json:"token"`
}

func addJxtv(key string, name string, id string){
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Referer: "http://www.jxntv.cn/",
		urlFmt: "https://live.jxtvcn.com.cn/live-jxtv/%s.m3u8?source=pc&t=%s&token=%s",
		urlBuildFunc: func(refererInfo RerferInfo) string {
			charset := []rune("ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678oOLl9gqVvUuI1")
			eTags := make([]rune,8)
			for  i := 0; i < 8; i ++ {
				eTags[i] = charset[rand.Intn(len(charset))]
			}
			timestamp := strconv.FormatInt(time.Now().Unix(),10)
			eTag := string(eTags)
			h := md5.New()
			h.Write([]byte(timestamp + id + ".m3u8"  + eTag + "dfasg2df!f"))
			a := h.Sum(nil)


			values := url.Values{}
			values.Set("t", timestamp)
			values.Set("stream", id + ".m3u8")

			req, err := http.NewRequest("POST", "https://app.jxntv.cn/Qiniu/liveauth/getPCAuth.php", strings.NewReader(values.Encode()))
			if err != nil {
				return ""
			}
			req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
			req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
			req.Header.Set("Origin", `http://www.jxntv.cn`)
			req.Header.Set("Referer", `http://www.jxntv.cn/`)
			req.Header.Set("Host", `app.jxntv.cn`)
			req.Header.Set("etag", eTag)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
			req.Header.Set("Content-Length", "34")
			req.Header.Set("Authorization", hex.EncodeToString(a))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return ""
			}

			result := JxtvAuthResult{}

			json.Unmarshal(body, &result)

			if result.Code != 200 {
				return ""
			}

			return fmt.Sprintf(refererInfo.urlFmt, refererInfo.Id, result.T, result.Token)
		},
	}
}

type gztvResult struct {
	Title       string `json:"title"`
	EntryType   string `json:"entry_type"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	PubDate     string `json:"pub_date"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Author      string `json:"author"`
	StreamUrl   string `json:"stream_url"`
}

func addGztv(key string, name string, id string){
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Referer: "https://www.gzstv.com/",
		urlFmt: "https://api.gzstv.com/v1/tv/%s/",
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := gztvResult{}

			json.Unmarshal(body, &result)

			http.RedirectHandler(result.StreamUrl,302).ServeHTTP(w, r)
		},
	}
}

func addCdtv(key string, name string, id string){
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Referer: "https://www.cditv.cn/",
		urlFmt: "https://www.cditv.cn/api.php?op=live&type=playTv&fluency=sd&videotype=m3u8&catid=192&id=%s",
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(string(body),302).ServeHTTP(w, r)
		},
	}
}

type AnhuiTvResult struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Aspect     string `json:"aspect"`
	CommentNum int    `json:"comment_num"`
	Snap       struct {
		Host     string `json:"host"`
		Dir      string `json:"dir"`
		Path     string `json:"path"`
		Filepath string `json:"filepath"`
		Filename string `json:"filename"`
	} `json:"snap"`
	SiteId     int    `json:"site_id"`
	ClickNum   int    `json:"click_num"`
	PraiseNum  int    `json:"praise_num"`
	ShareNum   int    `json:"share_num"`
	M3U8       string `json:"m3u8"`
	CurProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"cur_program"`
	NextProgram struct {
		StartTime string `json:"start_time"`
		Program   string `json:"program"`
	} `json:"next_program"`
	Logo struct {
		Square struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Path     string `json:"path"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"square"`
		Rectangle struct {
			Host     string `json:"host"`
			Dir      string `json:"dir"`
			Path     string `json:"path"`
			Filepath string `json:"filepath"`
			Filename string `json:"filename"`
		} `json:"rectangle"`
	} `json:"logo"`
	SaveTime      string `json:"save_time"`
	AudioOnly     string `json:"audio_only"`
	ContentUrl    string `json:"content_url"`
	ChannelStream []struct {
		Url        string `json:"url"`
		Name       string `json:"name"`
		StreamName string `json:"stream_name"`
		M3U8       string `json:"m3u8"`
		Bitrate    string `json:"bitrate"`
		StreamUrl  string `json:"stream_url"`
	} `json:"channel_stream"`
	NodeId   int    `json:"node_id"`
	ShareUrl string `json:"share_url"`
	OriginId int    `json:"origin_id"`
}
//        "appid":"m2otdjzyuuu8bcccnq",
//                    "appkey":"5eab6b4e1969a8f9aef459699f0d9000",
func addAhtv(key string, name string, id string){
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Referer: "http://www.ahtv.cn/",
		urlFmt: "http://mapi.ahtv.cn/api/open/ahtv/channel.php?appid=m2otdjzyuuu8bcccnq&appkey=5eab6b4e1969a8f9aef459699f0d9000&is_audio=0&category_id=1%2C2",
		urlBuildFunc: func(refererInfo RerferInfo) string {
			return refererInfo.urlFmt
		},
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := make([]AnhuiTvResult,0)

			_ =json.Unmarshal(body, &result)
			id,_:= strconv.Atoi(refererInfo.Id)
			for _, value := range result{
				if value.Id == id {
					w.Header().Set("Content-Type", "audio/x-mpegurl")
					http.RedirectHandler(string(value.M3U8),302).ServeHTTP(w, r)
					return
				}
			}
		},
	}
}

type HbTvResult struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Alias          string `json:"alias"`
	Stream         string `json:"stream"`
	Icon           string `json:"icon"`
	LiveType       string `json:"live_type"`
	AccessType     string `json:"access_type"`
	Fms            string `json:"fms"`
	Created        string `json:"created"`
	Createdby      string `json:"createdby"`
	PlaybillStatus string `json:"playbill_status"`
	Replay         string `json:"replay"`
	ReplayExpire   string `json:"replay_expire"`
	HasHd          string `json:"has_hd"`
	PublishedPc    string `json:"published_pc"`
	PublishedPhone string `json:"published_phone"`
	Rate           string `json:"rate"`
	DefaultThumb   string `json:"default_thumb"`
	StreamHd       string `json:"stream_hd"`
	Rtmp           string `json:"rtmp"`
	RtmpHd         string `json:"rtmp_hd"`
	Sort           string `json:"sort"`
	State          string `json:"state"`
	Url            string `json:"url"`
	PlayUrl        string `json:"play_url"`
	PlayUrlSd      string `json:"play_url_sd"`
	PlayUrlHd      string `json:"play_url_hd"`
	UsePub         string `json:"use_pub"`
	Shift          string `json:"shift"`
	ShiftStarttime string `json:"shift_starttime"`
	ShiftEndtime   string `json:"shift_endtime"`
	VirtualLive    string `json:"virtual_live"`
	PlayControl    string `json:"play_control"`
	VmsTid         string `json:"vms_tid"`
	ControlUrl     string `json:"control_url"`
	ShiftDay       string `json:"shift_day"`
	Token          string `json:"token"`
	Identity       string `json:"identity"`
	OssId          string `json:"oss_id"`
	OssHdId        string `json:"oss_hd_id"`
	Type           string `json:"type"`
}
func addHbtv(key string, name string, id string){
	referinfos[key] = &RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Referer: "http://app.cjyun.org/",
		urlFmt: "http://app.cjyun.org/video/player/stream?stream_id=%s&site_id=10008",
		afterFunc: func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				http.Error(w, http.StatusText(503), 503)
			}

			result := HbTvResult{}

			_ =json.Unmarshal(body, &result)

			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(string(result.Stream),302).ServeHTTP(w, r)
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
func init() {
	addCetv("CETV-1", "cetv1")
	addCetv("CETV-2", "cetv2")
	addCetv("CETV-3", "cetv3")
	addCetv("CETV-4", "cetv4")
	addGdtv("GDWS", "43")
	addGdtv("GDZJ", "44")
	addGdtv("GDXW", "45")
	addGdtv("GDTY", "47")
	addGdtv("NFWS", "51")
	addGdtv("GDJJKJ", "49")
	addGdtv("GDYS", "53")
	addGdtv("GDZY", "16")
	addGdtv("GDGJ", "46")
	addGdtv("GDSE", "54")
	addGdtv("GDJJKT", "66")
	addGdtv("GDNFGW", "42")
	addGdtv("GDLNXQ", "15")
	addGdtv("GDFC", "67")
	addGdtv("GDXDJY", "13")
	addGdtv("GDYD", "74")
	addGdtv("GRTNWHPD", "75")
	addSztv("SZWS", "AxeFRth")
	addSztv("SZYL", "1q4iPng")
	addSztv("SZSE", "1SIQj6s")
	addSztv("SZGG", "2q76Sw2")
	addSztv("SZDSJ", "4azbkoY")
	addSztv("SZDB", "9zoW71b")
	addSztv("SZYHGW", "BJ5u5k2")
	addSztv("SZDS", "ZwxzUXr")
	addSztv("SZGJ", "sztvgjpd")
	addSztv("SZGJ", "sztvgjpd")
	addSztv("SZTYJK", "sztvtyjk")
	addSztv("SZLG", "uGzbXhS")
	addSztv("SZYDSX", "wDF6KJ3")
	addSztv("SZDVSH", "xO1xQFv")
	addCqtv("CQWS", "重庆卫视 HD", "4918")
	addXjtv("XJSE", "新疆少儿 HD", "zb12")
	addXjtv("XJWS", "新疆卫视 HD", "zb01")
	addXjtv("XJWWE", "新疆维吾尔语综合 HD", "zb02")
	addXjtv("XJHSK", "新疆哈萨克语综合 HD", "zb03")
	addXjtv("XJHYZY", "新疆汉语综艺 HD", "zb04")
	addXjtv("XJWWEYS", "新疆维吾尔影视 HD", "zb05")
	addXjtv("XJHYJJSH", "新疆汉语经济生活 HD", "zb07")
	addXjtv("XJHSKZY", "新疆哈萨克语综艺 HD", "zb08")
	addXjtv("XJWWEJJSH", "新疆维吾尔经济生活 HD", "zb09")
	addXjtv("XJHYTYJK", "新疆汉语体育健康 HD", "zb10")
	addXjtv("XJHYXXFW", "新疆汉语信息服务 HD", "zb10")
	addJxtv("JXWS", "江西卫视 HD", "tv_jxtv1")
	addJxtv("JXDS", "江西都市频道 HD", "tv_jxtv2")
	addJxtv("JXJJSH", "江西经济生活频道 HD", "tv_jxtv3_hd")
	addJxtv("JXYSLY", "江西影视旅游频道 HD", "tv_jxtv4")
	addJxtv("JXGGNY", "江西公共农业频道 HD", "tv_jxtv5")
	addJxtv("JXSE", "江西少儿频道 HD", "tv_jxtv6")
	addJxtv("JXXW", "江西新闻频道 HD", "tv_jxtv7")
	addJxtv("JXYD", "江西移动电视 HD", "tv_jxtv8")
	addJxtv("JXFSGW", "江西风尚购物 HD", "tv_fsgw")
	addJxtv("JXTC", "江西陶瓷频道 HD", "tv_taoci")
	addGztv("GZWS", "贵州卫视 HD", "ch01")
	addGztv("GZGG", "贵州公共 HD", "ch02")
	addGztv("GZYSWY", "贵州影视文艺 HD", "ch03")
	addGztv("GZDZSH", "贵州大众生活 HD", "ch04")
	addGztv("GZD5", "贵州第5频道 HD", "ch05")
	addGztv("GZKJJK", "贵州科教健康 HD", "ch06")
	addGztv("GZSZYD", "贵州数字移动 HD", "ch07")
	addCdtv("CDXWZH", "成都新闻综合 HD", "1")
	addCdtv("CDJJZX", "成都经济咨询 HD", "2")
	addCdtv("CDDSSH", "成都都市生活 HD", "3")
	addCdtv("CDYSWY", "成都影视文艺 HD", "4")
	addCdtv("CDGG", "成都公共 HD", "5")
	addCdtv("CDSE", "成都少儿 HD", "6")
	addAhtv("AHWS", "安徽卫视 HD", "47")
	addAhtv("AHJJSH", "安徽经济生活 HD", "71")
	addAhtv("AHZYTY", "安徽综艺体育 HD", "73")
	addAhtv("AHYS", "安徽影视 HD", "72")
	addAhtv("AHGJ", "安徽国际 HD", "50")
	addAhtv("AHNYKJ", "安徽农业科教 HD", "51")
	addAhtv("AHGJ", "安徽国际 HD", "70")
	addAhtv("AHYD", "安徽移动电视 HD", "68")
	addAhtv("AHJC", "睛彩安徽 HD", "85")
	addHbtv("HBLS", "湖北垄上 HD", "438")
	addHbtv("MJGW", "湖北美嘉购物 HD", "439")
	addHbtv("HBJY", "湖北教育 HD", "437")
	addHbtv("HBSH", "湖北生活 HD", "436")
	addHbtv("HBYS", "湖北影视 HD", "435")
	addHbtv("HBGG", "湖北公共 HD", "434")
	addHbtv("HBZH", "湖北综合 HD", "433")
	addHbtv("HBJS", "湖北经视 HD", "432")
	addHbtv("HBWS", "湖北卫视 HD", "431")
	addHbtv("HBXW", "湖北新闻 HD", "470")
}

func AESDownloadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	src := r.Form.Get("src")
	keySrc := r.Form.Get("key")
	file := r.Form.Get("file")
	referer := r.Form.Get("referer")
	req, err := http.NewRequest("GET", src+"/"+file, nil)

	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
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

	key, _ := hex.DecodeString(keySrc)

	block, err := aes.NewCipher(key)


	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	origData := make([]byte, len(body))
	blockMode.CryptBlocks(origData, body)
	w.Header().Add("Content-Type", "video/MP2T")
	w.Write(origData)
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	rUrl := r.Form.Get("url")
	referer := r.Form.Get("referer")
	prefix := r.Form.Get("prefix")
	key := r.Form.Get("key")
	proxy := r.Form.Get("proxy")
	srcUrl := r.Form.Get("srcUrl")
	println("==" + prefix)
	//urlb, err := base64.StdEncoding.DecodeString(url)
	//refererb, err := base64.StdEncoding.DecodeString(referer)
	//prefixb, err := base64.StdEncoding.DecodeString(prefix)

	req, err := http.NewRequest("GET", rUrl, nil)

	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
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
	var newbody []byte
	if proxy != "" {
		suffix := ""
		if key != "" && srcUrl != "" {
			values := url.Values{}
			values.Set("src", srcUrl)
			values.Set("key", key)
			suffix = values.Encode()
		}

		if strings.HasSuffix(proxy, "/") {
			newbody = reRegx.ReplaceAll(body, []byte(proxy+"$0?"+suffix))
		} else if strings.HasSuffix(proxy, "=") {
			newbody = reRegx.ReplaceAll(body, []byte(proxy+"$0&"+suffix))
		} else {
			newbody = reRegx.ReplaceAll(body, []byte(proxy+"/$0?"+suffix))
		}
		w.Write(newbody)
	} else if prefix != "" {
		if strings.HasSuffix(prefix, "/") {

			newbody = reRegx.ReplaceAll(body, []byte(prefix+"$0"))
		} else {
			println(prefix + "/$0")
			newbody = reRegx.ReplaceAll(body, []byte(prefix+"/$0"))
		}
		w.Write(newbody)
	} else {
		w.Write(body)
	}

}

func RefererHandler(w http.ResponseWriter, r *http.Request) {
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
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if referinfo.beforeFunc != nil {
		referinfo.beforeFunc(*referinfo, url, req.Header)
	} else {
		if referinfo.Referer != "" {
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
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			var newbody []byte
			if strings.HasSuffix(referinfo.prefix, "/"){
				newbody = reRegx.ReplaceAll(body, []byte(referinfo.prefix+"$0"))
			} else {
				newbody = reRegx.ReplaceAll(body, []byte(referinfo.prefix+"/$0"))
			}


			logrus.Info(string(newbody))
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			w.Write(newbody)
		} else {
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			w.Write(body)
		}
	} else {
		referinfo.afterFunc(*referinfo, url, resp.Body, w, r)
	}

}		
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
		if referinfo.prefix != "" {
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			var newbody []byte
			if strings.HasSuffix(referinfo.prefix, "/") {
				newbody = reRegx.ReplaceAll(body, []byte(referinfo.prefix+"$0"))
			} else {
				newbody = reRegx.ReplaceAll(body, []byte(referinfo.prefix+"/$0"))
			}

			logrus.Info(string(newbody))
			w.Header().Set("Content-Type", "audio/x-mpegurl")

			w.Write(newbody)
		} else {
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			w.Write(body)
		}
	} else {
		referinfo.afterFunc(*referinfo, url, resp.Body, w, r)
	}

}
