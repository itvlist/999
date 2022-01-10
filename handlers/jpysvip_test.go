package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	sys_url "net/url"
	"strings"
	"testing"
	"time"
	"wmenjoy.com/iptv/utils"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func TestGetDetail5(t *testing.T)  {

	h := md5.New()
	h.Write([]byte("roZ68Okc5MUTMraMzb0161dc2453"))
	a := h.Sum(nil)
	println(hex.EncodeToString(a))

	block, err := aes.NewCipher([]byte("roZ68Okc5MUTMraM"))
	if err != nil {
		fmt.Println("")
	}
	data := "5bfc96eaef8d1e7ec0af912b9c8c88c194f54dcc3bc988aea89d322c20406457dab3ecd7bfb6145b92d8f945935f2deefb0f8bf018f194c15a5644d8f6758422e1ea17ef3f57e25ca13d97762dd9326ed06a1e8a732e170458a44239a82d1df6b5da76771dae969367221a5b2b96d6be44e6bf8b6f245760d55bbfa8aee242f116edc695e99d5a39824aca8eb56fe87af7087eab532cfd98169a003a3b071aaa694cafbd57c82437cfa7c2d35b113ba47a64f2e2bc7d67c4ec3bf1fd7d7d71c69eaaaf6ba6bf2532a0950c4f3dcce8ec73a79c1857ae85f10d9691e7ed51c52ccfcd4923f223d031efc2bb6a910f95758309d88271dcb9196d99eaaf2ad361ead6a9a6bce204a686b4afd70f2474c8cdb500201fb9c7b8ae60e0c091e1ed6cd6b931f904be51d0d22dc1e6f07b843f21e442df2a70c903622f989f20217e97c0412c330e0016f95724736c09cc9552f12d08d43ead4843db1adfa599ef2bcc3e359851b32b94e37991ca7d993de54f7cd94198cbd18fc3fb9a029537aa243f16f60e15988a441370d4019ae258db69f4e223d532a7cdd26eeeb11cfef1c4683a"
	result, _ := hex.DecodeString(data)
	//blockSize := block.BlockSize()
	//tmp := ZeroPadding([]byte(data), blockSize)
	mode := cipher.NewCBCDecrypter(block,[]byte("7384627385960726"))
	text :=  make([]byte, len(result))
	mode.CryptBlocks(text ,result)
	//
	fmt.Println(string(PKCS7UnPadding(text)))

}

func TestGetDetail4(t *testing.T){
	curr := time.Now()
	fmt.Println(curr.Format("2006"))
	fmt.Println(curr.Format("1"))
	fmt.Println(curr.Format("2"))

	str := fmt.Sprintf("chinamcloud%s年%s月%s日",curr.Format("2006"),curr.Format("1"),curr.Format("2"))

	block, err := aes.NewCipher([]byte("hzb%#02jn2*&23r#"))
	if err != nil {
		fmt.Println("")
	}
	//blockSize := block.BlockSize()
	//tmp := ZeroPadding([]byte(data), blockSize)
	data := PKCS7Padding([]byte(str), block.BlockSize())
	mode := cipher.NewCBCEncrypter(block,[]byte("1563432177954301"))
	ciphertext :=  make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)

	fmt.Println(strings.ToUpper(hex.EncodeToString(ciphertext)))
}
//MacPlayerConfig.player_list={"mx771":{"show":"MX\u84dd\u5149","des":"","ps":"1","parse":""},"anyiyun":{"show":"\u4e91\u6d77\u84dd\u5149","des":"","ps":"1","parse":""},"renrenmi":{"show":"RR\u84dd\u5149","des":"","ps":"1","parse":""},"rx":{"show":"rx","des":"rx","ps":"1","parse":""},"mengxin886":{"show":"mengxin886","des":"mengxin886","ps":"1","parse":""},"ltm3u8":{"show":"\u84dd\u5149\u7ebf\u8def","des":"","ps":"1","parse":"https:\/\/jxapp.jpysvip.net\/m3u8.php?url="},"dplayer":{"show":"\u84dd\u5149\u6781\u901f","des":"dplayer.js.org","ps":"1","parse":"\/dplayer\/analysis.php?v="},"niuxyun":{"show":"\u9996\u9009\u84dd\u5149","des":"","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"189pan":{"show":"\u84dd\u5149\u7ebf\u8def","des":"\u5fae\u76d8","ps":"1","parse":"\/dplayer\/analysis.php?v="},"mm":{"show":"PP\u4e91","des":"","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"fanqie":{"show":"\u756a\u8304\u8d44\u6e90","des":"fqzy.cc","ps":"1","parse":"https:\/\/jx.fqzy.cc\/jx.php?url="},"ddyunp":{"show":"\u591a\u591a\u84dd\u5149","des":"\u901f\u5ea6\u5feb\uff0c\u6e05\u6670\u5ea6\u9ad8\u3002","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"fuckapp":{"show":"\u95ea\u7535\u84dd\u5149","des":"\u84dd\u5149\u8d44\u6e90","ps":"1","parse":"\/dplayer\/analysis.php?v="},"xinm3u8":{"show":"\u4e0d\u5361\u7ebf\u8def","des":"\u65e0\u9700\u4e0b\u8f7d\u64ad\u653e\u5668","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"iframe":{"show":"\u8d85\u5feb\u4e91\u64ad","des":"iframe\u5916\u94fe\u6570\u636e","ps":"0","parse":""},"aliplayer":{"show":"\u963f\u91cc\u4e91","des":"","ps":"1","parse":"https:\/\/www.2ajx.com\/vip.php?url="},"funshion":{"show":"\u98ce\u884c\u89c6\u9891","des":"\u65e0\u9700\u5b89\u88c5\u4efb\u4f55\u63d2\u4ef6\uff0c\u9ad8\u901f\u64ad\u653e\u3002","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"videojs":{"show":"videojs-H5\u64ad\u653e\u5668","des":"videojs.com","ps":"0","parse":""},"iva":{"show":"H5\u81ea\u5e26\u89e3\u6790","des":"videojj.com","ps":"0","parse":""},"xigua":{"show":"\u897f\u74dc\u5f71\u97f3","des":"xigua.com","ps":"0","parse":""},"ffhd":{"show":"\u975e\u51e1\u5f71\u97f3","des":"www.feifan.com","ps":"0","parse":""},"wasu":{"show":"\u534e\u6570tv","des":"wasu.cn","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"letv":{"show":"\u4e50\u89c6\u89c6\u9891","des":"letv.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"mgtv":{"show":"\u8292\u679ctv","des":"mgtv.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"pptv":{"show":"pptv","des":"pptv.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"135m3u8":{"show":"\u6781\u901f\u4e91\u64ad","des":"","ps":"0","parse":""},"ckm3u8":{"show":"\u9177\u4e91","des":"","ps":"0","parse":""},"migu":{"show":"\u54aa\u5495\u89c6\u9891","des":"","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"youku":{"show":"\u4f18\u9177\u89c6\u9891","des":"youku.com","ps":"1","parse":"\/dplayer\/analysis.php?v="},"qq":{"show":"\u817e\u8baf\u89c6\u9891","des":"v.qq.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"sohu":{"show":"\u641c\u72d0\u89c6\u9891","des":"v.sohu.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"qiyi":{"show":"\u5947\u827a\u89c6\u9891","des":"qiyi.com","ps":"1","parse":"https:\/\/nbjx.vip\/?url="},"wjm3u8":{"show":"\u7ebf\u8def\u2460","des":"","ps":"1","parse":"https:\/\/jx.xhswglobal.com\/dplayer\/?url="},"dbm3u8":{"show":"\u767e\u5ea6\u4e91M3U8","des":"\u5728\u7ebf\u64ad\u653e","ps":"1","parse":"\/dplayer\/analysis.php?v="},"zuidam3u8":{"show":"\u9ad8\u6e05\u4e91","des":"\u6700\u5927","ps":"1","parse":"https:\/\/www.jpysvip.net\/jhjx\/?url="},"kuyun":{"show":"kuyun","des":"https:\/\/www.mahuazy.com\/","ps":"0","parse":""},"mahua":{"show":"\u7ebf\u8def1","des":"","ps":"1","parse":"\/dplayer\/analysis.php?v="},"subom3u8":{"show":"\u6781\u901f\u4e91\u64ad","des":"","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"bjm3u8":{"show":"\u6781\u901f\u84dd\u5149","des":"\u53ea\u80fd\u624b\u673a\u64ad\u653e\u3002","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"},"ltnb":{"show":"\u9ad8\u901f\u84dd\u5149","des":"","ps":"1","parse":"https:\/\/jxapp.jpysvip.net\/m3u8.php?url="},"appplayer":{"show":"\u79d2\u64ad\u4e0d\u5931\u6548","des":"","ps":"1","parse":"https:\/\/www.jpysvip.net\/appplayer.html?"}},MacPlayerConfig.downer_list={"http":{"show":"http\u4e0b\u8f7d","des":"des\u63d0\u793a\u4fe1\u606f","ps":"1","parse":""},"xunlei":{"show":"\u8fc5\u96f7\u4e0b\u8f7d","des":"des\u63d0\u793a\u4fe1\u606f","ps":"1","parse":""}},MacPlayerConfig.server_list={"server1":{"show":"\u6d4b\u8bd5\u670d\u52a1\u56681","des":"des\u63d0\u793a\u4fe1\u606f1"}};
func TestGetDetail(t *testing.T){
	detail, err := getDetailInfo("https://www.jpysvip.net/voddetail/144499.html")
	assert.Nil(t, err)

	print(detail.Title)
}

func TestGetVodDetail(t *testing.T){
	detail, err := getVodPlayDetailInfo("https://www.jpysvip.net/vodplay/144499-1-1.html")
	assert.Nil(t, err)

	print(detail.Url)
}


func TestSign(t *testing.T) {
	url := fmt.Sprintf("https://nbjx.vip/?url=%s", "https://v.qq.com/x/cover/mzc00200ihwjf82/d0034ptm2h9.html")

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
	print("\n")

	print(realResult.Url)
	print("\n")

}