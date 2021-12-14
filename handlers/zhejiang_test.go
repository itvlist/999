package handlers

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"testing"
)
const str =`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:186291246
#EXT-X-TARGETDURATION:8
#EXT-X-PROGRAM-DATE-TIME:2020-07-20T09:12:34+0800
#EXTINF:8.00
1080p/0/000082931ed3d499_4d005e8_30d60.ts
#EXTINF:8.00
1080p/0/000082931ed3d499_4d31348_30be8.ts
#EXTINF:8.00
1080p/0/000082931ed3d499_4d61f30_30d60.ts`

func TestReplace(t *testing.T) {

	re := regexp.MustCompilePOSIX(`([^#]+\.ts)`)
	logrus.Printf(string(re.ReplaceAll([]byte(str), []byte("/test/$0"))))
}
