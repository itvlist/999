package priovices

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"wmenjoy.com/iptv/handlers"
)
//玛纳斯综合频道 http://218.84.127.245:1026/hls/main1/playlist.m3u8
func init()  {
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
	handlers.RefererInfos[key] = &handlers.RerferInfo{
		Id:     id,
		Key:    key,
		Jump:   true,
		Name:   name,
		Prefix: "http://livehyw5.chinamcache.com/hyw/",
		UrlFmt: "http://livehyw5.chinamcache.com/hyw/%s.m3u8?txSecret=%s&txTime=%s",
		UrlBuildFunc: func(refererInfo handlers.RerferInfo) string {
			return "http://mediaxjtvs.chinamcache.com/hyw/media/playerJson/liveChannel/7d40edeb62fe4f8a9d9a08bc653dcab6_PlayerParamProfile.json"
		},
		AfterFunc: func(refererInfo handlers.RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request) {
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
			realUrl := fmt.Sprintf(refererInfo.UrlFmt, refererInfo.Id, hex.EncodeToString(a), timestamp)
			w.Header().Set("Content-Type", "audio/x-mpegurl")
			http.RedirectHandler(realUrl, 302).ServeHTTP(w, r)
		},
	}
}

