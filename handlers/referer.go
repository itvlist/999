package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"wmenjoy.com/iptv/utils"
)

type RerferInfo struct {
	Id           string
	Key          string
	Name         string
	Group        string
	Category     string
	Quality      string
	Referer      string
	UrlFmt       string
	Jump         bool
	DirectReturn bool
	ReRegxp      *regexp.Regexp
	Prefix       string
	UrlBuildFunc func(refererInfo RerferInfo) string
	BeforeFunc   func(refererInfo RerferInfo, url string, header http.Header)
	AfterFunc    func(refererInfo RerferInfo, url string, bodyReader io.ReadCloser, w http.ResponseWriter, r *http.Request)
}

var RefererInfos = make(map[string]*RerferInfo, 0)

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
	prefix := r.Form.Get("Prefix")
	key := r.Form.Get("key")
	proxy := r.Form.Get("proxy")
	srcUrl := r.Form.Get("srcUrl")
	//urlb, err := base64.StdEncoding.DecodeString(url)
	//refererb, err := base64.StdEncoding.DecodeString(referer)
	//prefixb, err := base64.StdEncoding.DecodeString(Prefix)

	req, err := http.NewRequest("GET", rUrl, nil)

	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
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
	urlstr := ""

	id := r.Form.Get("id")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	referinfo := RefererInfos[strings.ToUpper(id)]

	if referinfo == nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	if referinfo.UrlBuildFunc == nil {
		urlstr = fmt.Sprintf(referinfo.UrlFmt, referinfo.Id)
	} else {
		urlstr = referinfo.UrlBuildFunc(*referinfo)
	}

	if referinfo.DirectReturn {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl;charset=UTF-8")
		http.RedirectHandler(urlstr, 302).ServeHTTP(w, r)
	}

	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if referinfo.BeforeFunc != nil {
		referinfo.BeforeFunc(*referinfo, urlstr, req.Header)
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

	if resp.StatusCode == 302 {
		location, _ := resp.Location()

		values := url.Values{}
		values.Set("url", location.String())
		values.Set("prefix", location.String()[0:strings.LastIndex(location.String(), "/")])

		urlstring := fmt.Sprintf("http://%s:8880/transfer?%s", utils.GetIP(), values.Encode())
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl;charset=UTF-8")
		http.RedirectHandler(urlstring, 302).ServeHTTP(w, r)
	}

	if referinfo.AfterFunc == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 503)
			return
		}
		if referinfo.Prefix != "" {
			prefix := referinfo.Prefix

			if strings.Contains(prefix, "%s") {
				prefix = fmt.Sprintf(prefix, referinfo.Id)
			}

			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl;charset=UTF-8")
			var newbody []byte
			regexc := referinfo.ReRegxp
			if regexc == nil {
				regexc = reRegx
			}

			if strings.HasSuffix(prefix, "/") {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"$0"))
			} else {
				newbody = regexc.ReplaceAll(body, []byte(prefix+"/$0"))
			}
			w.Write(newbody)
		} else {
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl;charset=UTF-8")
			w.Write(body)
		}
	} else {
		referinfo.AfterFunc(*referinfo, urlstr, resp.Body, w, r)
	}

}
