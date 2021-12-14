package main

import (
	"flag"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"net/http"
	"wmenjoy.com/iptv/handlers"
)


func main() {
	flag.IntVar(&handlers.DefaultThreadNum, "t", 5, "Multi Download Thread Num")
	flag.Parse()
	logrus.SetOutput(colorable.NewColorableStdout())
	http.HandleFunc("/index.m3u8", indexHandler)
	//http.HandleFunc("/api/", apiHandler)
	//http.HandleFunc("/videoplayback/", videoHandler)
	http.HandleFunc("/sitv.m3u8", handlers.SitvHandler)
	http.HandleFunc("/byr.m3u8", handlers.ByrHandler)
	http.HandleFunc("/youtube.m3u8", handlers.YoutubeIndexHandler)
	http.HandleFunc("/litv.m3u8", handlers.LitvHandler)
	http.HandleFunc("/4gtv.m3u8", handlers.FourgtvHandler)
	http.HandleFunc("/fjtv.m3u8", handlers.FjtvHandler)
	http.HandleFunc("/grtn.m3u8", handlers.GrtnHandler)
	http.HandleFunc("/youjia.m3u8", handlers.YoujiaHandler)
	http.HandleFunc("/inews.m3u8", handlers.InewsHandler)
	http.HandleFunc("/douyu.m3u8",handlers.DouyuHandler)
	http.HandleFunc("/neu6.m3u8", handlers.Neu6Handler)
	http.HandleFunc("/neu.m3u8", handlers.NeuHandler)
	http.HandleFunc("/tuna.m3u8", handlers.TunaHandler)
	http.HandleFunc("/live/pool/", handlers.FourgtvTsHandler)
	http.HandleFunc("/live/", handlers.QmHandler)
	http.HandleFunc("/hls/", handlers.Neu6tsHandler)
	http.HandleFunc("/neu/hls/", handlers.NeutsHandler)
	http.HandleFunc("/tuna/hls/", handlers.TunaTsHandler)
	http.HandleFunc("/sdlive/", handlers.SctvHandler)
	http.HandleFunc("/hdlive/", handlers.SctvHandler)
	http.HandleFunc("/haixia/", handlers.FjtvApiHandler)
	http.HandleFunc("/haixia_sd/", handlers.FjtvApiHandler)
	http.HandleFunc("/youjia/", handlers.YoujiaApiHandler)
	http.HandleFunc("/api/", handlers.YoutubeApiHandler)
	http.HandleFunc("/videoplayback/", handlers.YoutubeVideoHandler)
	http.HandleFunc("/hi/vod/", handlers.FourgtvApiHandler)
	http.HandleFunc("/zj/api", handlers.ZhejiangApiHandler)
	http.HandleFunc("/zj/", handlers.ZhejiangHandler)
	http.HandleFunc("/", handlers.ByrApiHandler)
	http.ListenAndServe(":8080", nil)

}