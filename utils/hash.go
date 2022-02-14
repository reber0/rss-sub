/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-02-14 14:37:10
 * @LastEditTime: 2022-02-14 15:01:09
 */
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// base64 encode
func B64Encode(plainText []byte) string {
	return base64.StdEncoding.EncodeToString(plainText)
}

// base64 decode
func B64Decode(cipherText string) []byte {
	plainText, _ := base64.StdEncoding.DecodeString(cipherText)
	return plainText
}

// md5 加密
func Md5(plainText []byte) string {
	return fmt.Sprintf("%x", md5.Sum(plainText))
}

// Sha1 加密
func Sha1(plainText []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(plainText))
}

// AES padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// AES unpadding
func PKCS7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}

// AES 加密, CBC
func AesEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = PKCS7Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

// AES 解密
func AesDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = PKCS7UnPadding(plainText)
	return plainText, nil
}
