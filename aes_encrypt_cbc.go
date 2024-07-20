package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesEncryptCBC(origData []byte, k []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)              // 补全码
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                   // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                // 加密
	return encrypted
}
