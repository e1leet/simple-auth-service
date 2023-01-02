package manager

import (
	"crypto/sha256"
	"encoding/hex"
)

type Manager interface {
	CheckPassword(password, hashedPassword string) bool
	HashPassword(password string) string
}

type manager struct {
	salt string
}

func New(salt string) Manager {
	return &manager{
		salt: salt,
	}
}

func (m *manager) CheckPassword(password, hashedPassword string) bool {
	return m.HashPassword(password) == hashedPassword
}

func (m *manager) HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password + m.salt))

	return hex.EncodeToString(h.Sum(nil))
}
