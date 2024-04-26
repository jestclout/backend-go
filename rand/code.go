package rand

import (
	"bytes"
	"math/rand"
)

const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// NewGameCode generate and returns a token that can be used in a url.
func NewGameCode(length int) (string, error) {
	buf := bytes.NewBuffer(nil)

	for i := 0; i < length; i++ {
		ch := validChars[rand.Intn(len(validChars))]
		err := buf.WriteByte(ch)
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
