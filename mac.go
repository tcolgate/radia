package vonq

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
)

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	log.Println("AGH", string(message), messageMAC, expectedMAC)
	return hmac.Equal(messageMAC, expectedMAC)
}

// GenMAC
func GenMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
