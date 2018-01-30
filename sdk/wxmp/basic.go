package wxmp

import (
	"sort"
	"crypto/sha1"
	"strings"
	"fmt"
)

// Common message header
type MessageHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

func CheckSignature(token string, timestamp string, nonce string, signature string) bool {
	strs := []string{token, timestamp, nonce}
	sort.Strings(strs) //order
	orderedStr := strings.Join(strs,"") //combine

	//sha1 hash
	h := sha1.New()
	h.Write([]byte(orderedStr))
	sha1Value := fmt.Sprintf("%x", h.Sum(nil))

	return sha1Value == signature
}

