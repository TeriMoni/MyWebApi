package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(parm string) string {
	h := md5.New()
	h.Write([]byte(parm)) // 需要加密的字符串
	return hex.EncodeToString(h.Sum(nil))

}
