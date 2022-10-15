package util

import (
	"fmt"
	"net/http"
	"os"

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

func ValidateFileType(dst string) (string, error) {
	file, err := os.Open(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}
