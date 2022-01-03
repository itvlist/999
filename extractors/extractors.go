package extractors

import (
	"net/url"
	"strings"
	"wmenjoy.com/iptv/extractors/acfun"
	"wmenjoy.com/iptv/extractors/bcy"
	"wmenjoy.com/iptv/extractors/bilibili"
	"wmenjoy.com/iptv/extractors/douyin"
	"wmenjoy.com/iptv/extractors/douyu"
	"wmenjoy.com/iptv/extractors/eporner"
	"wmenjoy.com/iptv/extractors/facebook"
	"wmenjoy.com/iptv/extractors/geekbang"
	"wmenjoy.com/iptv/extractors/haokan"
	"wmenjoy.com/iptv/extractors/instagram"
	"wmenjoy.com/iptv/extractors/iqiyi"
	"wmenjoy.com/iptv/extractors/mgtv"
	"wmenjoy.com/iptv/extractors/miaopai"
	"wmenjoy.com/iptv/extractors/netease"
	"wmenjoy.com/iptv/extractors/pixivision"
	"wmenjoy.com/iptv/extractors/pornhub"
	"wmenjoy.com/iptv/extractors/qq"
	"wmenjoy.com/iptv/extractors/streamtape"
	"wmenjoy.com/iptv/extractors/tangdou"
	"wmenjoy.com/iptv/extractors/tiktok"
	"wmenjoy.com/iptv/extractors/tumblr"
	"wmenjoy.com/iptv/extractors/twitter"
	"wmenjoy.com/iptv/extractors/types"
	"wmenjoy.com/iptv/extractors/udn"
	"wmenjoy.com/iptv/extractors/universal"
	"wmenjoy.com/iptv/extractors/vimeo"
	"wmenjoy.com/iptv/extractors/weibo"
	"wmenjoy.com/iptv/extractors/xvideos"
	"wmenjoy.com/iptv/extractors/yinyuetai"
	"wmenjoy.com/iptv/extractors/youku"
	"wmenjoy.com/iptv/extractors/youtube"
	"wmenjoy.com/iptv/utils"
)

var extractorMap map[string]types.Extractor

func init() {
	douyinExtractor := douyin.New()
	youtubeExtractor := youtube.New()
	stExtractor := streamtape.New()

	extractorMap = map[string]types.Extractor{
		"": universal.New(), // universal extractor

		"douyin":     douyinExtractor,
		"iesdouyin":  douyinExtractor,
		"bilibili":   bilibili.New(),
		"bcy":        bcy.New(),
		"pixivision": pixivision.New(),
		"youku":      youku.New(),
		"youtube":    youtubeExtractor,
		"youtu":      youtubeExtractor, // youtu.be
		"iqiyi":      iqiyi.New(iqiyi.SiteTypeIqiyi),
		"iq":         iqiyi.New(iqiyi.SiteTypeIQ),
		"mgtv":       mgtv.New(),
		"tangdou":    tangdou.New(),
		"tumblr":     tumblr.New(),
		"vimeo":      vimeo.New(),
		"facebook":   facebook.New(),
		"douyu":      douyu.New(),
		"miaopai":    miaopai.New(),
		"163":        netease.New(),
		"weibo":      weibo.New(),
		"instagram":  instagram.New(),
		"twitter":    twitter.New(),
		"qq":         qq.New(),
		"yinyuetai":  yinyuetai.New(),
		"geekbang":   geekbang.New(),
		"pornhub":    pornhub.New(),
		"xvideos":    xvideos.New(),
		"udn":        udn.New(),
		"tiktok":     tiktok.New(),
		"haokan":     haokan.New(),
		"acfun":      acfun.New(),
		"eporner":    eporner.New(),
		"streamtape": stExtractor,
		"streamta":   stExtractor, // streamta.pe
	}
}

// Extract is the main function to extract the data.
func Extract(u string, option types.Options) ([]*types.Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, err
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
