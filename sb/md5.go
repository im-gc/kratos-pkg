package sb

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(bs []byte) (string, error) {

	h := md5.New()
	if _, err := h.Write(bs); nil != err {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
