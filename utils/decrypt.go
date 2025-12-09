package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"jtyl_bitable/global"
	"strings"
)

func Decrypt(encrypt string) string {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return ""
	}
	if len(buf) < aes.BlockSize {
		return ""
	}
	keyBs := sha256.Sum256([]byte(global.CONFIG.Event.EncryptKey))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return ""
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	if len(buf)%aes.BlockSize != 0 {
		return ""
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1])
}
