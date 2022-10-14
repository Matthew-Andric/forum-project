package util

import (
	"fmt"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var (
	S securecookie.SecureCookie
)

func SaltPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}
