package udpping

import (
	"crypto/hmac"
	"crypto/sha256"
)

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

// GenMAC
func genMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
