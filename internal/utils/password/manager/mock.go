package manager

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) CheckPassword(password, hashedPassword string) bool {
	args := m.Called(password, hashedPassword)

	return args.Bool(0)
}

func (m *Mock) HashPassword(password string) string {
	args := m.Called(password)

	return args.String(0)
}
