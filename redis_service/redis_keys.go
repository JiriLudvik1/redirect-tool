package redis_service

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

const redirectPrefix = "redirect_"

func createUrlHash(url string) string {
	hash := md5.New()
	hash.Write([]byte(url))
	hashedBytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedBytes)
}

func getRedirectRedisKey(hash string) string {
	return fmt.Sprintf("%s%s", redirectPrefix, hash)
}
