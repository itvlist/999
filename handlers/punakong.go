package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func PunaKongHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Request URL:%s", r.URL)
	url := ""

	url = fmt.Sprintf("https://%s", r.URL.Path[len("/punakong/"):])
	logrus.Info(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
	req.Header.Set("Accept", `*/*`)
	req.Header.Set("Referer", `https://www.zxzj.fun/`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	/*url = fmt.Sprintf("http://hw-m-l.cztv.com/%s", r.URL.Path[4:])*/
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	bodyStr := string(body)

	reg1 := regexp.MustCompile(`var url = '[^']+'`)
	result := reg1.FindString(bodyStr)
	if len(result) > 0{
		url = result[len("var url = '"):len(result) -1 ]
		url = decodeStr(htoStr(strReverse(url)))
		http.RedirectHandler(url,302).ServeHTTP(w, r)
	} else {
		http.Error(w, err.Error(), 503)
		return
	}
}

func decodeStr(str string)  string{
	i := (len(str) - 6)/2
	return str[0:i] + str[i+6:]
}

func strReverse(str string) string {
	strLen := len(str)

	result := ""
	for i := 0; i < strLen; i++ {
		result += string(str[strLen - i - 1 ])
	}
	return result
}


func htoStr(str string) string {
	index := 0
	result := ""
	for ; index < len(str); index = index + 2 {
		value := str[index:index +2]
		numvalue,_ := strconv.ParseInt(value, 0x10, 32)
		result += string(rune(numvalue))
	}
	return result
}

