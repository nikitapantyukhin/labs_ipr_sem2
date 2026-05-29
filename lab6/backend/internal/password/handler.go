package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const separator = "@"

type IPasswordHandler interface {
	HashPassword(password string) []byte
	VerifyPassword(rawPassword string, encodedPassword []byte) (bool, error)
}

type Handler struct {
	cfg *PasswordConfig
}

func (h Handler) generateSalt() []byte {
	bytes := make([]byte, h.cfg.SaltSize)
	_, _ = rand.Read(bytes)
	return bytes
}

func (h Handler) HashPassword(password string) []byte {
	salt := h.generateSalt()
	hashedPassword := h.hashPassword(password, salt)
	return []byte(
		fmt.Sprintf(
			"%s%s%s",
			base64.RawStdEncoding.EncodeToString(salt),
			separator,
			base64.RawStdEncoding.EncodeToString(hashedPassword),
		),
	)
}

func (h Handler) hashPassword(password string, salt []byte) []byte {
	hashedPassword := argon2.IDKey(
		[]byte(password),
		salt,
		h.cfg.Iterations,
		h.cfg.Memory,
		h.cfg.Threads,
		h.cfg.KeyLength,
	)

	return hashedPassword
}

func (h Handler) decodePasswordAndSalt(encodedPassword []byte) (salt []byte, password []byte, err error) {
	str := string(encodedPassword)
	if !strings.Contains(str, separator) {
		return nil, nil, errors.New("can't decode password as its of wrong format")
	}

	splittedString := strings.Split(str, separator)
	if len(splittedString) != 2 {
		return nil, nil, errors.New("can't decode password as its of wrong format")
	}

	b64Salt, b64Password := splittedString[0], splittedString[1]

	salt, saltDecodeError := base64.RawStdEncoding.DecodeString(b64Salt)
	if saltDecodeError != nil {
		return nil, nil, saltDecodeError
	}

	password, passwordDecodeError := base64.RawStdEncoding.DecodeString(b64Password)
	if passwordDecodeError != nil {
		return nil, nil, passwordDecodeError
	}

	return
}

func (h Handler) VerifyPassword(rawPassword string, encodedPassword []byte) (bool, error) {
	salt, hashedPassword, passwordFormatError := h.decodePasswordAndSalt(encodedPassword)
	if passwordFormatError != nil {
		return false, passwordFormatError
	}

	otherHashedPassword := h.hashPassword(rawPassword, salt)
	return subtle.ConstantTimeCompare(hashedPassword, otherHashedPassword) == 1, nil

}

func CreateHandler(cfg *PasswordConfig) IPasswordHandler {
	return &Handler{cfg: cfg}
}
