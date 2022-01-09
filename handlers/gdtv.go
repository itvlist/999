package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func _gererateHeader(key string, url string, timestamp string) string  {
	secret := []byte("dfkcY1c3sfuw0Cii9DWjOUO3iQy2hqlDxyvDXd1oVMxwYAJSgeB6phO8eW1dfuwX")
	message := []byte("GET\nhttps://gdtv-api.gdtv.cn/api/tv/v2/tvChannel/54?tvChannelPk=54&node=NjEyMDk3N2U4YzM4NjhjY2UxZDZlMDcxNjk4Mzc3ZmUtbkZ4azFyZUxBODRqQlI2VmJwc09VTUZjJTJCdnY0M2E1cjVZeWkyeTM0eldrdjZub2NBNmJkVVVKbWtVVDRKWjByNkxSZ0EwdTJjSlZvdU9GMU1Vc3ZhU2ZicFBSeWs4clI5bGkydWc4WlNVZyUzRA==\n1641700205023\n")

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}