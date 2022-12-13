package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

// 验证签名
func CheckSignature(signature, timestamp, nonce, token string) bool {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	sum := sha1.Sum([]byte(sl[0] + sl[1] + sl[2]))
	return signature == hex.EncodeToString(sum[:])
}
