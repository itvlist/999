package provinces

import (
	"regexp"
	"wmenjoy.com/iptv/utils"
)

var host = utils.GetIP()

var reRegx2 = regexp.MustCompilePOSIX(`([^#]+\.ts)`)

