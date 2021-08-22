package security

import (
	"crypto/md5"
)

func MD5(src []byte) []byte {
	hash := md5.New()
	hash.Write(src)
	return hash.Sum(nil)
}
