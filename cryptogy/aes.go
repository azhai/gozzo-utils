package cryptogy

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"strings"
)

type PaddingFunc func(cipherText []byte, blockSize int) []byte
type UnpaddingFunc func(origData []byte) []byte

func ZeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(cipherText, padtext...)
}

func ZeroUnpadding(origData []byte) []byte {
	return bytes.TrimFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize //需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(cipherText, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padtext...)
}

func PKCS7Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES加密，支持模式CBC、CFB、CTR、OFB，不支持ECB和GCM
// 其中CBC模式一般需要填充，用法: c.SetPaddingFunc("PKCS5")
type AESCipher struct {
	modeName  string
	iv        []byte
	Padding   PaddingFunc
	Unpadding UnpaddingFunc
	cipher.Block
}

func NewAESCipher(mode string, key []byte) (*AESCipher, error) {
	c := &AESCipher{modeName: strings.ToUpper(mode)}
	block, err := aes.NewCipher(key)
	if err == nil {
		c.Block = block
		c.iv = key[:c.BlockSize()]
	}
	return c, err
}

func (c *AESCipher) SetPaddingFunc(name string) {
	switch strings.ToUpper(name) {
	case "0", "ZERO":
		c.Padding = ZeroPadding
		c.Unpadding = ZeroUnpadding
	case "5", "PKCS5":
		c.Padding = PKCS5Padding
		c.Unpadding = PKCS5Unpadding
	case "7", "PKCS7":
		c.Padding = PKCS7Padding
		c.Unpadding = PKCS7Unpadding
	}
}

func (c *AESCipher) GetStream(isDecrypt bool) cipher.Stream {
	switch c.modeName {
	default:
		return nil
	case "CFB":
		if isDecrypt {
			return cipher.NewCFBDecrypter(c.Block, c.iv)
		} else {
			return cipher.NewCFBEncrypter(c.Block, c.iv)
		}
	case "CTR":
		return cipher.NewCTR(c.Block, c.iv)
	case "OFB":
		return cipher.NewOFB(c.Block, c.iv)
	}
}

func (c *AESCipher) GetEncrypter() cipher.BlockMode {
	return cipher.NewCBCEncrypter(c.Block, c.iv)
}

func (c *AESCipher) GetDecrypter() cipher.BlockMode {
	return cipher.NewCBCDecrypter(c.Block, c.iv)
}

func (c *AESCipher) Encrypt(origData []byte) ([]byte, error) {
	if c.Padding != nil {
		origData = c.Padding(origData, c.BlockSize())
	}
	cipherText := make([]byte, len(origData))
	if c.modeName == "CBC" {
		c.GetEncrypter().CryptBlocks(cipherText, origData)
	} else {
		c.GetStream(false).XORKeyStream(cipherText, origData)
	}
	return cipherText, nil
}

func (c *AESCipher) Decrypt(cipherText []byte) ([]byte, error) {
	origData := make([]byte, len(cipherText))
	if c.modeName == "CBC" {
		c.GetDecrypter().CryptBlocks(origData, cipherText)
	} else {
		c.GetStream(true).XORKeyStream(cipherText, origData)
	}
	if c.Unpadding != nil {
		origData = c.Unpadding(origData)
	}
	return origData, nil
}
