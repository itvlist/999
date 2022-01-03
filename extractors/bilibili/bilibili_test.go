package bilibili

import (
	"testing"

	"wmenjoy.com/iptv/extractors/types"
	"wmenjoy.com/iptv/test"
)

func TestBilibili(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "normal test 1",
			args: test.Args{
				URL:   "https://www.bilibili.com/video/av20203945/",
				Title: "【2018拜年祭单品】相遇day by day",
			},
			playlist: false,
		},
		{
			name: "normal test 2",
			args: test.Args{
				URL:   "https://www.bilibili.com/video/av41301960",
				Title: "【英雄联盟】2019赛季CG 《觉醒》",
				Size:  70696896,
			},
			playlist: false,
		},
		{
			name: "bangumi test",
			args: test.Args{
				URL:   "https://www.bilibili.com/bangumi/play/ep167000",
				Title: "狐妖小红娘：第70话 苏苏智商上线",
			},
		},
		{
			name: "bangumi playlist test",
			args: test.Args{
				URL:   "https://www.bilibili.com/bangumi/play/ss5050",
				Title: "一人之下：第1话 异人刀兵起，道炁携阴阳",
			},
			playlist: true,
		},
		{
			name: "playlist test",
			args: test.Args{
				URL:   "https://www.bilibili.com/video/av16907446/",
				Title: "\"不要相信歌词，他们为了押韵什么都干得出来\"",
			},
			playlist: true,
		},
		{
			name: "bangumi movie test",
			args: test.Args{
				URL:   "https://www.bilibili.com/bangumi/play/ss12044",
				Title: "你的名字。",
			},
		},
		{
			name: "8k test",
			args: test.Args{
				URL:   "https://www.bilibili.com/video/BV1qM4y1w716",
				Title: "【8K演示片】B站首发！你的设备还顶得住吗？",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				data []*types.Data
				err  error
			)
			if tt.playlist {
				// for playlist, we don't check the data
				_, err = New().Extract(tt.args.URL, types.Options{
					Playlist:     true,
					ThreadNumber: 9,
				})
				test.CheckError(t, err)
			} else {
				data, err = New().Extract(tt.args.URL, types.Options{})
				test.CheckError(t, err)
				test.Check(t, tt.args, data[0])
			}
		})
	}
}