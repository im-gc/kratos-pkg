package sb

import (
	"encoding/json"
)

func MarshalAndMd5(data interface{}) (string, string, error) {

	bs, err := json.Marshal(data)
	if nil != err {
		return "", "", err
	}

	md5, err := Md5(bs)
	if nil != err {
		return "", "", err
	}
	return B2S(bs), md5, nil
}

func UnMarshal[T any](s *string, v T) (T, error) {

	bs := S2B(s)
	if err := json.Unmarshal(bs, &v); nil != err {
		return v, err
	}
	return v, nil
}
