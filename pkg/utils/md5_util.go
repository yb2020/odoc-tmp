package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 计算字符串的MD5哈希值
func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
