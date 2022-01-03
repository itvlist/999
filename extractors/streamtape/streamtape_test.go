package streamtape

import (
	"testing"

	"wmenjoy.com/iptv/extractors/types"
	"wmenjoy.com/iptv/test"
)

func TestStreamtape(t *testing.T) {
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://streamtape.com/e/vkoKlwYPo9F4mRo",
				Title: "annie.mp4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.args.URL, types.Options{})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
	}
}
