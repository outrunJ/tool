package tool

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func MD5(text []byte) []byte {
	hash := md5.New()
	hash.Write(text)
	return hash.Sum(nil)
}

func MD5String(text string) string {
	return hex.EncodeToString(MD5([]byte(text)))
}

func MD5RString(text []byte) string {
	return hex.EncodeToString(MD5(text))
}

func Base64Encode(origData []byte) string {
	return base64.StdEncoding.EncodeToString(origData)
}
func Base64Decode(crypted string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(crypted)
}

func SHA1(text []byte) []byte {
	sb := sha1.Sum(text)
	return sb[:]
}

func SHA1RString(text []byte) string {
	b := sha1.Sum(text)
	return hex.EncodeToString(b[:])
}

func Encrypt(origData, key []byte) (string, error) {
	bits, err := DesEncrypt(origData, key)
	if err != nil {
		return "", err
	}
	return Base64Encode(bits), nil
}

func Decrypt(crypted string, key []byte) ([]byte, error) {
	bits, err := Base64Decode(crypted)
	if err != nil {
		return nil, err
	}
	return DesDecrypt(bits, key)
}

func Sign(origData []byte, salt []byte) string {
	return Base64Encode(SHA1(append(origData, salt...)))
}

func UniqueID(b []byte) string {
	id := MD5RString(b) + SHA1RString(b)
	return id
}
