package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return RandomCode() // Bad idea maybe?
	}
	return fmt.Sprintf("%06d", n)
}
