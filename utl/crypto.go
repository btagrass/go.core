package utl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

var (
	Key = []byte("0123456789ABCDEFFEDCBA9876543210")
	iv  = []byte("0123456776543210")
)

// 解密
func Decrypt(data string) (string, error) {
	dataSrc, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	cbc := cipher.NewCBCDecrypter(block, iv)
	dataDst := make([]byte, len(dataSrc))
	cbc.CryptBlocks(dataDst, dataSrc)
	dataDst = unPaddingPkcs5(dataDst)

	return string(dataDst), nil
}

// 摘要
func Digest(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	bytes := md.Sum(nil)

	return base64.StdEncoding.EncodeToString(bytes)
}

// 加密
func Encrypt(data string) (string, error) {
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	dataSrc := paddingPkcs5([]byte(data), block.BlockSize())
	cbc := cipher.NewCBCEncrypter(block, iv)
	dataDst := make([]byte, len(dataSrc))
	cbc.CryptBlocks(dataDst, dataSrc)

	return base64.StdEncoding.EncodeToString(dataDst), nil
}

// HmacSha256
func HmacSha256(data string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(data))
	bytes := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(bytes)
}

// Md5
func Md5(data string) string {
	bytes := md5.Sum([]byte(data))

	return fmt.Sprintf("%x", bytes)
}

// 填充PKCS5
func paddingPkcs5(data []byte, blockSize int) []byte {
	paddingSize := blockSize - len(data)%blockSize
	paddings := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	return append(data, paddings...)
}

// 取消填充PKCS5
func unPaddingPkcs5(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
