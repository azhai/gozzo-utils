package cryptogy

import (
	"crypto/hmac"
	"encoding/base64"
	"hash"
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

func (h MacHash) SetKey(key string) MacHash {
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
