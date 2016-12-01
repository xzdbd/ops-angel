package models

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

const (
	TOKEN = "opsangel"
)

func CheckSignature(timestamp string, nonce string) string {
	strs := sort.StringSlice{TOKEN, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
