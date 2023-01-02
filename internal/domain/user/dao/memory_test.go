package dao

import (
	"context"
	"io"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.Output(io.Discard)
}

func TestMemoryDAO_Create(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		dao := NewMemory()
		id, err := dao.Create(context.Background(), &UserStorage{
			Username: "test",
			Password: "test",
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})
	t.Run("username_already_used", func(t *testing.T) {
		dao := NewMemory()

		const username = "test"

		_, err := dao.Create(context.Background(), &UserStorage{
			Username: username,
			Password: "test",
		})
		assert.NoError(t, err)

		_, err = dao.Create(context.Background(), &UserStorage{
			Username: username,
			Password: "test",
		})
		assert.ErrorIs(t, err, ErrUsernameAlreadyUsed)
	})
}

func TestMemoryDAO_GetById(t *testing.T) {
	t.Run("get_existing", func(t *testing.T) {
		dao := NewMemory()
		id, err := dao.Create(context.Background(), &UserStorage{
			Username: "test",
			Password: "test",
		})
		assert.NoError(t, err)

		usr, err := dao.GetById(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, usr.ID)
	})

	t.Run("non_existent", func(t *testing.T) {
		dao := NewMemory()
		_, err := dao.GetById(context.Background(), 123)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestMemoryDAO_GetByUsername(t *testing.T) {
	t.Run("get_existing", func(t *testing.T) {
		dao := NewMemory()
		expected := &UserStorage{
			Username: "test",
			Password: "test",
		}
		id, err := dao.Create(context.Background(), expected)
		assert.NoError(t, err)

		usr, err := dao.GetByUsername(context.Background(), expected.Username)
		assert.NoError(t, err)
		assert.Equal(t, id, usr.ID)
	})

	t.Run("not_existent", func(t *testing.T) {
		dao := NewMemory()
		_, err := dao.GetByUsername(context.Background(), "test")
		assert.ErrorIs(t, err, ErrUserNotFound)
	})
}

func TestMemoryDAO_DeleteById(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dao := NewMemory()
		expected := &UserStorage{
			Username: "test",
			Password: "test",
		}
		id, err := dao.Create(context.Background(), expected)
		assert.NoError(t, err)

		err = dao.DeleteById(context.Background(), id)
		assert.NoError(t, err)

		_, err = dao.GetById(context.Background(), id)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("non-existent", func(t *testing.T) {
		dao := NewMemory()
		err := dao.DeleteById(context.Background(), 123)
		assert.NoError(t, err)
	})
}
