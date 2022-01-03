package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	sys_url "net/url"
	"strings"
	"testing"
	"wmenjoy.com/iptv/utils"
)

func TestSign(t *testing.T) {
	url := fmt.Sprintf("https://nbjx.vip/?url=%s", "https://v.qq.com/x/cover/im60meg91bo9dbr.html")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
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
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	strBody := string(body)
	time := utils.MatchOneOf(strBody, `var\s+time\s*=\s*'([^']+)'`)[1]
	wap := utils.MatchOneOf(strBody, `var\s+wap\s*=\s*'([^']+)'`)[1]
	realUrl := utils.MatchOneOf(strBody, `var\s+url\s*=\s*'([^']+)'`)[1]
	vkey := utils.MatchOneOf(strBody, `var\s+vkey\s*=\s*'([^']+)'`)[1]
	fvkey := utils.MatchOneOf(strBody, `var\s+fvkey\s*=\s*'([^']+)'`)[1]

	url = "https://nbjx.vip/xmflv-1.SVG"
	//dataParam :=fmt.Sprintf("time=%s&wap=%s&url=%s&vkey=%s&fvkey=%s", time,wap,realUrl, vkey, Sign(fvkey))

	dataParam := sys_url.Values{}
	dataParam.Set("time", time)
	dataParam.Set("wap", wap)
	dataParam.Set("url", realUrl)
	dataParam.Set("vkey", vkey)
	dataParam.Set("fvkey", Sign(fvkey))




	print(dataParam)
	req, err = http.NewRequest("POST", url,  strings.NewReader(dataParam.Encode()))
	if err != nil {
		return
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
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var realResult jpysvip

	err = json.Unmarshal(body, &realResult)
	if err != nil {
		print(string(body))
		return
	}
	print(string(body))

	print(realResult.Url)
}