package String

import (
	"crypto/md5"
	"encoding/hex"
)

// 将字符串转换成md5
func ToMd5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
