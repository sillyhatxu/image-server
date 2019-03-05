package aes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	src := make([]string, 31)
	for i := 1; i <= 31; i++ {
		var msg string
		if i < 10 {
			msg = "0" + strconv.Itoa(i)
		} else {
			msg = strconv.Itoa(i)
		}
		encrypted, err := Encrypt(msg)
		fmt.Printf("encrypted : %v\n", encrypted)
		assert.Nil(t, err)
		src[i-1] = encrypted
	}
	for i := 1; i <= 31; i++ {
		var msg string
		if i < 10 {
			msg = "0" + strconv.Itoa(i)
		} else {
			msg = strconv.Itoa(i)
		}
		decrypted, err := Decrypt(src[i-1])
		fmt.Printf("decrypted : %v\n", decrypted)
		assert.Nil(t, err)
		assert.Equal(t, msg, decrypted)
	}
}
