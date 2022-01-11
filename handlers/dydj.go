package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wmenjoy.com/iptv/utils"
)

type DiDyListItem struct {
	Name string
	SimpleName string
	Group string
	Num   int
	Image   string
	Url    string
	StartSeason int
	EndSeason int
	MaxEp    int
}

type DiDyDetail struct {
	Type         string `json:"type"`
	Tracklist    bool   `json:"tracklist"`
	Tracknumbers bool   `json:"tracknumbers"`
	Images       bool   `json:"images"`
	Artists      bool   `json:"artists"`
	Tracks       []struct {
		Src         string `json:"src"`
		Src0        string `json:"src0"`
		Src1        string `json:"src1"`
		Src2        string `json:"src2"`
		Src3        string `json:"src3"`
		Title       string `json:"title"`
		Type        string `json:"type"`
		Caption     string `json:"caption"`
		Description string `json:"description"`
		Image       struct {
			Src    string `json:"src"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"image"`
		Thumb struct {
			Src    string `json:"src"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumb"`
		Meta struct {
			LengthFormatted string `json:"length_formatted"`
		} `json:"meta"`
		Portn      string      `json:"portn"`
		Srctype    string      `json:"srctype"`
		Cut        string      `json:"cut"`
		Vttshift   string      `json:"vttshift"`
		UserIP     interface{} `json:"userIP"`
		Subsrc     string      `json:"subsrc"`
		Dimensions struct {
			Original struct {
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"original"`
			Resized struct {
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"resized"`
		} `json:"dimensions"`
	} `json:"tracks"`
}

func didyEncrypt(path string) string{
	etimes := time.Now().Unix() + 600

	uTxt := "{\"path\":\"" + path + "\",\"expire\":" + strconv.FormatInt(etimes * 1000, 10) + "}"
	block, err := aes.NewCipher([]byte("zevS%th@*8YWUm%K"))
	if err != nil {
		return ""
	}
	mode := cipher.NewCBCEncrypter(block,[]byte("5080305495198718"))
	data := PKCS5Padding([]byte(uTxt), block.BlockSize())
	ciphertext :=  make([]byte, len(data))

	mode.CryptBlocks(ciphertext, data)

	return url.QueryEscape(base64.StdEncoding.EncodeToString(ciphertext))
}

func getAlbumInfo(name string)(simpleName string, maxEp int, startSeason int, endSeason int){

	index := strings.LastIndex(name, "(")
	simpleName = name
	if index != -1 {
		simpleName = name[0:index]
	}
	index = strings.LastIndex(simpleName, "第")
	if  index != -1 {
		simpleName = simpleName[0:index]
	}

	simpleName = strings.TrimSpace(simpleName)

	if strings.Contains(name, "季"){
		result := utils.MatchOneOf(name, "([0-9]+)-([0-9])+")

		if len(result) == 3 {
			endSeason, _ = strconv.Atoi(result[2])
			startSeason, _ =strconv.Atoi(result[1])

			if strings.Contains(name, "+"){
				return simpleName, -1,  startSeason, endSeason + 1
			}


			return simpleName, -1, startSeason, endSeason
		}
		result = utils.MatchOneOf(name, "([0-9]+).*季.*([0-9])+")
		if len(result) == 3 {
			endSeason, _ = strconv.Atoi(result[1])
			maxEp, _= strconv.Atoi(result[2])
			return simpleName, maxEp, 1, endSeason
		}

		result = utils.MatchOneOf(name, "更新.*([0-9]+).*季")
		if len(result) == 2 {
			endSeason, _ = strconv.Atoi(result[1])
			return simpleName,-1, 1, endSeason
		}
		result = utils.MatchOneOf(name, "季.*([0-9])+")
		if len(result) == 2 {
			maxEp, _= strconv.Atoi(result[1])
			return simpleName, maxEp, -1, -1
		}

		return simpleName,-1, -1, -1

	} else {

		result := utils.MatchOneOf(name, "[0-9]+")

		if len(result) == 1 {
			maxEp, _= strconv.Atoi(result[0])
			return simpleName, maxEp, -1, -1
		}

		return simpleName,-1, -1 ,-1
	}
}

type DiDyDetailInfo struct{
	Season int
	Items []DiDyDetailItem

}
type DiDyDetailItem struct {
	Name string
	SimpleName string
	Group string
	Image   string
	Url    string
}

func getDidyDetail(name string, url string) *DiDyDetailInfo{
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept", `*/*`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	detailJson := utils.MatchOneOf(string(body), "<script[^>]+application/json\">([^<]+)</script>")[1]
	result := DiDyDetail{}

	_ =json.Unmarshal([]byte(detailJson), &result)

	detail := DiDyDetailInfo{
		Items: make([]DiDyDetailItem, 0),
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil
	}

	group := ""

	 groupItem := dom.Find(".meta_categories a")

	 if groupItem.Size() != 0 {
		 group = groupItem.Get(0).FirstChild.Data
	 }


	endSeason := dom.Find(".page-links .post-page-numbers").Size()
	for index, track  := range result.Tracks{
		detail.Season = endSeason
		detail.Items = append(detail.Items, DiDyDetailItem{
			Name: strconv.Itoa(index),
			Group: group,
			SimpleName: name,
			Url: strings.ReplaceAll(track.Src0,`\`, ""),

		})
	}

	return &detail
}
func getDidyList(url string) []DiDyListItem {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept", `*/*`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}

	itemList := make([]DiDyListItem, 0)
	dom.Find(".post-box-list article").Each(func(i int, selection *goquery.Selection) {
		group := selection.Find(".post-box-meta a").Get(0).FirstChild.Data
		item := selection.Find(".post-box-text a")
		imageText, exist := selection.Find(".post-box-image").Attr("style")

		if exist{
			imageText = utils.MatchOneOf(imageText, "(http[^)]+)")[0]
		}

		address, exist := selection.Attr("data-href")

		name := item.Get(item.Size() - 1).FirstChild.Data

		simpleName, maxEp, startSeason, EndSeason := getAlbumInfo(name)

		if exist{
			itemList = append(itemList, DiDyListItem{
				Group: group,
				Name: name,
				Url: address,
				Num: 1,
				Image: imageText,
				MaxEp: maxEp,
				StartSeason: startSeason,
				EndSeason: EndSeason,
				SimpleName: simpleName,
			})
		}

	})
	return itemList
}


func DiDyHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	url := ""

	id := r.Form.Get("id")
	ep := r.Form.Get("ep")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}

	if strings.HasPrefix(id, "https://ddys.tv/"){
		id = id[len("https://ddys.tv/"):]
	}

	if strings.HasSuffix(id, "/") {
		id = id[0:len(id)-1]
	}

	url = fmt.Sprintf("https://ddys.tv/%s/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept", `*/*`)

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
	detailJson := utils.MatchOneOf(string(body), "<script[^>]+application/json\">([^<]+)</script>")[1]
	result := DiDyDetail{}

	_ =json.Unmarshal([]byte(detailJson), &result)

	var src string
	if ep == ""{
		src = result.Tracks[0].Src0
	} else {
		index, _ := strconv.Atoi(ep)
		src = result.Tracks[index].Src0
	}

	w.Header().Set("Content-Type", "audio/x-mpegurl")

	realUrl := getDidyRealUrl(strings.ReplaceAll(src,`\`, ""))

	if realUrl == ""{
		http.Error(w, err.Error(), 503)
		return
	}

	http.RedirectHandler(realUrl, 302).ServeHTTP(w, r)

}

type DidyParserResult struct {
	Cache bool   `json:"cache"`
	Url   string `json:"url"`
}
func getDidyRealUrl(url string) string {
	url = fmt.Sprintf("https://v.ddys.tv:19543/video?type=mix&id=%s", didyEncrypt(url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept", `*/*`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}


	result := DidyParserResult{}

	_ =json.Unmarshal(body, &result)

	return result.Url

}
