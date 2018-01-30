package wxmp

import (
	"encoding/base64"
)

type EcryptedMsg struct {
	Encrypt string
	MsgSignature string
	TimeStamp string
	Nonce string
}

//decode encoding AES key
func decodeEncodingAesKey(key string) string{
	decodedKey, _ := base64.StdEncoding.DecodeString(key+"=")
	if len(decodedKey) != 32 {
		return ""
	}

	return string(decodedKey)
}

// DecryptMessage: 检验消息的真实性，并且获取解密后的明文
// @param sMsgSignature: 签名串，对应URL参数的msg_signature
// @param sTimeStamp: 时间戳，对应URL参数的timestamp
// @param sNonce: 随机串，对应URL参数的nonce
// @param sPostData: 密文，对应POST请求的数据
// xml_content: 解密后的原文，当return返回0时有效
// @return: 成功0，失败返回对应的错误码
// 验证安全签名
func DecryptMessage(msg string) {

}
