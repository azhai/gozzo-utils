package cryptogy

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

/*
// 可通过openssl产生
// openssl genrsa -out rsa_private_key.pem 4096
var privKey = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var pubKey = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtRs0cDbOgZwe7zgJPo11
lQyHrRrLDcO00LS2FH9NhdlQ8F/Wh6T6HUuuQhjVoKST3nqx3RE9mkzomw8tHdnA
0AU0THXLPptDshsclxZJZXzo9HZuqaPHWIFuWw7grCSGSjnls7Xakdc8ndPs3jnI
HKLf3douQjB0zmpj+60uuQ+48tG3R7skbSPIO5k7ji1HUTlhRAXzAiT5lhXTtMWx
3yEO17Dlu8x4ZXoWHd5fA3m3OsaonLmzOshW7uVMa4Ecb0ZgTdn9O9nR/dvZ8v+c
MhCTGkq/HgPuYgVCsqAdGBo9aUpamG/LB407mWKz4amjjIgM1DxU3hOv3J0oc6YU
pggXu0keev0xFwcWM36pVgMKEFwMsm4Bu7ZZ/AJSWsmPz0SC2QSsaRWuo7spRCGz
4CpGl2YvvGCQqytzgjB8BTn2dbKMBjUJjXtyw7FRXRy1WGgICBzayP+lPXRjYT6v
Hb+zy4oLV7vzwAk8fKB+mbDhv62Onb7vkOtLg7uiTvO5AHIGBYsxo92MIQG/gQop
0vJHRDlQKCgdmJbcrVWD4Jm5zWlS99PlUgMxRAYTS7JEv2/X9GDowdEHZY8/mUeJ
9dI8B6Iz65XGYKGLAAbUgXRoXFXw4V1dFD+dCCsseC5QKm+/6HfTIUwWfoEp40T2
UqWHmhg2z2Itx8B0m9cORIsCAwEAAQ==
-----END PUBLIC KEY-----
`
*/

// RSA加密，支持模式PKCS#1 v1.5，不支持OAEP/PSS
// 用法: NewRSACipher(privKey, pubKey).Sign(crypto.SHA256, []byte("hello world"))
type RSACipher struct {
	privKey, pubKey []byte
}

func NewRSACipher(privKey, pubKey string) RSACipher {
	return RSACipher{privKey: []byte(privKey), pubKey: []byte(pubKey)}
}

// 解密pem格式的公钥/私钥
func (c RSACipher) GetBlock(key []byte, errmsg string) (*pem.Block, error) {
	var err error
	block, _ := pem.Decode(key)
	if block == nil {
		err = errors.New(errmsg)
	}
	return block, err
}

func (c RSACipher) GetPublicKey() (*rsa.PublicKey, error) {
	block, err := c.GetBlock(c.pubKey, "public key error !")
	if err != nil {
		return nil, err
	}
	// 解析公钥
	face, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	return face.(*rsa.PublicKey), nil
}

func (c RSACipher) GetPrivateKey() (*rsa.PrivateKey, error) {
	block, err := c.GetBlock(c.privKey, "private key error!")
	if err != nil {
		return nil, err
	}
	// 解析PKCS1格式的私钥
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 加密
func (c RSACipher) Encrypt(origData []byte) ([]byte, error) {
	key, err := c.GetPublicKey()
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, key, origData)
}

// 解密
func (c RSACipher) Decrypt(cipherText []byte) ([]byte, error) {
	key, err := c.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, key, cipherText)
}

// 签名
func (c RSACipher) Sign(hash crypto.Hash, msg []byte) ([]byte, error) {
	key, err := c.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed[:])
}

// 校验
func (c RSACipher) Verify(hash crypto.Hash, msg, sig []byte) error {
	key, err := c.GetPublicKey()
	if err != nil {
		return err
	}
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(key, hash, hashed[:], sig)
}
