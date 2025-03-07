package conv

import (
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func StringToInt64(str string) (int64, error) {
	newData, err := strconv.ParseInt(str, 10, 64)

	return newData, err
}

func StringToInt(s string) (int, error) {
	numb, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return numb, err
}
