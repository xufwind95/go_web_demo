package auth

import "golang.org/x/crypto/bcrypt"

// 用bcrypt库对纯文本文件进行加密(go语言自带加密库)
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
