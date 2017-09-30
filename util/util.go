package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

var base64coder = base64.StdEncoding

func HashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Hash(s []byte) string {
	return fmt.Sprintf("%x", md5.Sum(s))
}
