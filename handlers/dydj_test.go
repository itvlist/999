package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestEncrypt(t *testing.T) {

	println(didyEncrypt("/v/movie/The.Lost.Daughter.mp4"))
}

func TestItemList(t *testing.T){
	firstList := getDidyList("https://ddys.tv/category/drama/western-drama/")

	for _, item := range firstList{

		detail := getDidyDetail(item.SimpleName, item.Url)

		for _, item := range detail.Items{
			fmt.Printf("%s,%s,%s\n", item.Group, item.SimpleName, item.Url)
		}

		if detail.Season > 1 {
			for i := 2; i <= detail.Season; i ++ {
				detail := getDidyDetail(item.SimpleName, item.Url + "/" + strconv.Itoa(i))
				for _, item := range detail.Items{
					fmt.Printf("%s,%s,%s\n", item.Group, item.SimpleName, item.Url)
				}
			}
		}
	}

}



func TestRegexItem(t *testing.T){

	name, maxEp, startSeason, endSeason := getAlbumInfo("毒枭 第1-3季")
	assert.Equal(t, "毒枭", name)
	assert.Equal(t, -1, maxEp)
	assert.Equal(t, 1, startSeason)
	assert.Equal(t, 3, endSeason)
	name, maxEp, startSeason, endSeason = getAlbumInfo("黑镜 第1-5季 (更新 第5季 Youtube特别篇)")
	assert.Equal(t, "黑镜", name)
	assert.Equal(t, -1, maxEp)
	assert.Equal(t, 1, startSeason)
	assert.Equal(t, 5, endSeason)
	name, maxEp, startSeason, endSeason = getAlbumInfo("新闻编辑室 第一季")
	assert.Equal(t, "新闻编辑室", name)

	assert.Equal(t, -1, maxEp)
	assert.Equal(t, -1, startSeason)
	assert.Equal(t, -1, endSeason)
	name, maxEp, startSeason, endSeason = getAlbumInfo("欧美剧热映中嗜血法医 杀魔新生 (更新至10)")
	assert.Equal(t, "欧美剧热映中嗜血法医 杀魔新生", name)

	assert.Equal(t, 10, maxEp)
	assert.Equal(t, -1, startSeason)
	assert.Equal(t, -1, endSeason)
	name, maxEp, startSeason, endSeason = getAlbumInfo("末日巡逻队 (第一季)")
	assert.Equal(t, "末日巡逻队", name)

	assert.Equal(t, -1, maxEp)
	assert.Equal(t, -1, startSeason)
	assert.Equal(t, -1, endSeason)
	name, maxEp, startSeason, endSeason = getAlbumInfo("毒枭: 墨西哥 (更新 第3季)")
	assert.Equal(t, "毒枭: 墨西哥", name)

	assert.Equal(t, -1, maxEp)
	assert.Equal(t, 1, startSeason)
	assert.Equal(t, 3, endSeason)
}