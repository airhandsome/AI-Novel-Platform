package utils

import "time"

func GenerateRandomString(length int) string {
	return time.Now().String()[:length]
}
