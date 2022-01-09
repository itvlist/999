package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	sys_url "net/url"
	"strings"
	"wmenjoy.com/iptv/utils"
)


type jpysvip struct {
	Player string `json:"player"`
	Success int `json:"success"`
	Type   string `json:"type"`
	Url    string `json:"url"`
}

type jpysvipPlayItem struct {
	SubTitle string
	PageUrl string
	PlayUrl string
}
type jpysvipPlayList struct {
	Id string
	SrcName string
	List []jpysvipPlayItem

}
type jpysvipDetail struct {
	Title string
	list []jpysvipPlayList
}

type jpysvipVodPlayDetail struct {
	Flag string `json:"flag"`
	Encrypt int `json:"encrypt"`
	Trysee int `json:"trysee"`
	Points int `json:"points"`
	Link string `json:"link"`
	LinkNext string `json:"link_next"`
	LinkPre string `json:"link_pre"`
	Url string `json:"url"`
	UrlNext string `json:"url_next"`
	From string `json:"from"`
	Server string `json:"server"`
	Note string `json:"note"`
	Id string `json:"id"`
	Sid int `json:"sid"`
	Nid int `json:"nid"`
}

type validKey struct {
	Time  string
	Wap   string
	Url   string
	Cip   string
	Vkey string
	Fvkey string
}


func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize//需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)//生成填充的文本
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)//用0去填充
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}


func Sign(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	a := h.Sum(nil)
	c :=hex.EncodeToString(a)

	block, err := aes.NewCipher([]byte(c))
	if err != nil {
		return ""
	}
	//blockSize := block.BlockSize()
	//tmp := ZeroPadding([]byte(data), blockSize)
	mode := cipher.NewCBCEncrypter(block,[]byte("ren163com5201314"))
	ciphertext :=  make([]byte, len([]byte(data)))
	mode.CryptBlocks(ciphertext, []byte(data))
	return base64.StdEncoding.EncodeToString(ciphertext)

}

func getDetailInfo(url string)(*jpysvipDetail, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authority", "www.jpysvip.net")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("sec-ch-ua-mobile","?0")
	req.Header.Set("sec-ch-ua-platform","macOS")
	req.Header.Set("upgrade-insecure-requests","1")
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site","none")
	req.Header.Set("sec-fetch-mode","navigate")
	req.Header.Set("sec-fetch-user","?1")
	req.Header.Set("sec-fetch-dest","document")
	req.Header.Set("referer","https://www.jpysvip.net/")
	req.Header.Set("accept-language","zh-CN,zh;q=0.9")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)

	detail := &jpysvipDetail{}
	detail.Title = utils.MatchOneOf(strBody, `class="title">([^<]+)</h1>`)[1]
	detail.list = make([]jpysvipPlayList, 0)
	matchResult := utils.MatchAll(strBody,"<a href=\"#playlist([^\"]+)\" [^>]+>([^<]+)" )
	for _,  value := range matchResult{
		if  value[2] == "极速蓝光" || value[2] == "极速云播" {
			continue
		}
		listId := value[1]

		playList := jpysvipPlayList{
			Id: listId,
			SrcName: value[2],
			List: make([]jpysvipPlayItem, 0),
		}

		matchResult2 := utils.MatchAll(strBody,"href=\"(/vodplay/[0-9]+-"+ listId+"-[0-9]+.html)\">([^<]+)" )
		for _,  value2 := range matchResult2 {
			playList.List = append(playList.List, jpysvipPlayItem{
				PageUrl: fmt.Sprintf("https://www.jpysvip.net%s", value2[1]),
				SubTitle: value2[2],
			})
		}

		detail.list = append(detail.list, playList)



	}

	return detail, nil
}

func getVodPlayDetailInfo(url string)(*jpysvipVodPlayDetail,error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authority", "www.jpysvip.net")
	req.Header.Set("upgrade-insecure-requests","1")
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site","none")
	req.Header.Set("sec-fetch-mode","navigate")
	req.Header.Set("sec-fetch-user","?1")
	req.Header.Set("sec-fetch-dest","document")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("sec-ch-ua-mobile","?0")
	req.Header.Set("sec-ch-ua-platform","macOS")
	req.Header.Set("referer","https://www.jpysvip.net/")
	req.Header.Set("accept-language","zh-CN,zh;q=0.9")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := utils.MatchOneOf(string(body), "player_aaaa=({[^}]+})")[1]

	bodyStr = strings.ReplaceAll(bodyStr,"\\","")
	var result jpysvipVodPlayDetail
	json.Unmarshal([]byte(bodyStr), &result)

	return &result, nil
}

func getValidKeyInfo(urlParam string)(*validKey, error){
	url := fmt.Sprintf("https://nbjx.vip/?url=%s", urlParam)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("sec-ch-ua-mobile","?0")
	req.Header.Set("sec-ch-ua-platform","macOS")
	req.Header.Set("upgrade-insecure-requests","1")
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site","cross-site")
	req.Header.Set("sec-fetch-mode","navigate")
	req.Header.Set("sec-fetch-dest","iframe")
	req.Header.Set("referer","https://www.jpysvip.net/")
	req.Header.Set("accept-language","zh-CN,zh;q=0.9")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)

	return &validKey{
		Time : utils.MatchOneOf(strBody, `var\s+time\s*=\s*'([^']+)'`)[1],
		Wap:  utils.MatchOneOf(strBody, `var\s+wap\s*=\s*'([^']+)'`)[1],
		Url : utils.MatchOneOf(strBody, `var\s+url\s*=\s*'([^']+)'`)[1],
		Cip : utils.MatchOneOf(strBody, `var\s+cip\s*=\s*'([^']+)'`)[1],
		Vkey : utils.MatchOneOf(strBody, `var\s+vkey\s*=\s*'([^']+)'`)[1],
		Fvkey : utils.MatchOneOf(strBody, `var\s+fvkey\s*=\s*'([^']+)'`)[1],
	}, nil

}

func getRealUrlForVipx(key validKey) (*jpysvip, error){
	url := "https://nbjx.vip/xmflv-1.SVG"
	dataParam := sys_url.Values{}
	dataParam.Set("time", key.Time)
	dataParam.Set("wap", key.Wap)
	dataParam.Set("url", key.Url)
	dataParam.Set("vkey", key.Vkey)
	dataParam.Set("fvkey", Sign(key.Fvkey))

	req, err := http.NewRequest("POST", url,  strings.NewReader(dataParam.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("authority","nbjx.vip")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("accept","application/json, text/javascript, */*; q=0.01")
	req.Header.Set("content-type","application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("x-requested-with","XMLHttpRequest")
	req.Header.Set("sec-ch-ua-mobile","?0")
	req.Header.Set("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform","macOS")
	req.Header.Set("origin","https://nbjx.vip")
	req.Header.Set("sec-fetch-site","same-origin")
	req.Header.Set("sec-fetch-mode","cors")
	req.Header.Set("sec-fetch-dest","empty")
	req.Header.Set("accept-language","zh-CN,zh;q=0.9")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var realResult jpysvip

	err = json.Unmarshal(body, &realResult)
	if err != nil {
		print(string(body))
		return nil,err
	}
	return &realResult, nil
}

func JspyVipHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	urlParam := r.Form.Get("url")
	if urlParam == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	keyInfo, err := getValidKeyInfo(urlParam)
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	jpysvipResult, err := getRealUrlForVipx(*keyInfo)


	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	http.RedirectHandler(jpysvipResult.Url,302).ServeHTTP(w, r)

}
