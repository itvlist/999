package handlers

import (
	"io/ioutil"
	"net/http"
)
//https://jhpc.manduhu.com/j1217.php?url=
//https://fendou.duoduozy.com/m3u8/202201/4/9aefd6f2c21b618657ffc402d10068f23592ed1f.m3u8?st=bRyYq1p_EkTlTcWXwc46pg&e=1641609068
func getDetailInfo2(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("authority", "fendou.duoduozy.com")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("accept", "*/*")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("Referer", "https://dp.duoduozy.com")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	strBody := string(body)

	print(strBody)
	return nil
}
