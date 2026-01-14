package util

import (
	"os"
	"strings"
)

func VerifyUser(userName string, password string) bool {
	if strings.ToLower(userName) != "affan" && strings.ToLower(userName) != "tiara" {
		return false
	}

	actualPassword := os.Getenv("USER_PASSWORD")
	if password != actualPassword {
		return false
	}

	return true
}
