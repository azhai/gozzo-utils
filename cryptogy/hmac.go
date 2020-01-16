package cryptogy

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"strings"
)

type NewHashFunc func() hash.Hash

// hmac哈希，例如 NewMacHash(sha256.New).SetKey("nonce")
type MacHash struct {
	creator   NewHashFunc
	secretKey []byte
}

func NewMacHash(creator NewHashFunc) *MacHash {
	return &MacHash{creator: creator}
}

func (h *MacHash) SetKey(key string) *MacHash {
	h.secretKey = []byte(key)
	return h
}

func (h MacHash) MacSum(text string) []byte {
	mac := hmac.New(h.creator, h.secretKey)
	mac.Write([]byte(text))
	return mac.Sum(nil)
}

func (h MacHash) Sign(text string) string {
	return base64.StdEncoding.EncodeToString(h.MacSum(text))
}

func (h MacHash) Verify(text, hashed string) bool {
	decoded, err := base64.StdEncoding.DecodeString(hashed)
	if err == nil {
		return hmac.Equal(decoded, h.MacSum(text))
	}
	return false
}

// 产生随机salt
func RandSalt(size int) string {
	buf := make([]byte, (size+1)/2)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)[:size]
	}
	return ""
}

// 带salt值的sha256密码哈希
type SaltPassword struct {
	saltLen int
	saltSep string
	*MacHash
}

func NewSaltPassword(len int, sep string) *SaltPassword {
	return &SaltPassword{
		saltLen: len, saltSep: sep,
		MacHash: NewMacHash(sha256.New),
	}
}

// 设置密码
func (p *SaltPassword) CreatePassword(plainText string) string {
	saltValue := RandSalt(p.saltLen)
	cipherText := p.SetKey(saltValue).Sign(plainText)
	return saltValue + p.saltSep + cipherText
}

// 校验密码
func (p *SaltPassword) VerifyPassword(plainText, cipherText string) bool {
	pieces := strings.SplitN(cipherText, p.saltSep, 2)
	if len(pieces) == 2 {
		return p.SetKey(pieces[0]).Verify(plainText, pieces[1])
	}
	return false
}
