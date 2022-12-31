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
		id, err := dao.Create(context.Background(), &SessionStorage{
			ExpiresIn: 123,
			UserID:    1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})
}

func TestMemoryDAO_GetByID(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dao := NewMemory()
		id, err := dao.Create(context.Background(), &SessionStorage{
			ExpiresIn: 123,
			UserID:    1,
		})
		assert.NoError(t, err)

		session, err := dao.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, session.ID)
	})

	t.Run("non_existent", func(t *testing.T) {
		dao := NewMemory()
		_, err := dao.GetByID(context.Background(), 123)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

func TestMemoryDAO_GetByToken(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dao := NewMemory()
		actual := &SessionStorage{
			ExpiresIn: 123,
			UserID:    1,
		}
		id, err := dao.Create(context.Background(), actual)
		assert.NoError(t, err)

		expected, err := dao.GetByToken(context.Background(), actual.Token)
		assert.NoError(t, err)
		assert.Equal(t, id, expected.ID)
	})

	t.Run("non_existent", func(t *testing.T) {
		dao := NewMemory()
		_, err := dao.GetByToken(context.Background(), "test")
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

func TestMemoryDAO_DeleteByID(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dao := NewMemory()
		id, err := dao.Create(context.Background(), &SessionStorage{
			ExpiresIn: 123,
			UserID:    1,
		})
		assert.NoError(t, err)

		err = dao.DeleteByID(context.Background(), id)
		assert.NoError(t, err)

		_, err = dao.GetByID(context.Background(), id)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

func TestMemoryDAO_DeleteByToken(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dao := NewMemory()
		session := &SessionStorage{
			ExpiresIn: 123,
			UserID:    1,
		}
		id, err := dao.Create(context.Background(), session)
		assert.NoError(t, err)

		err = dao.DeleteByToken(context.Background(), session.Token)
		assert.NoError(t, err)

		_, err = dao.GetByID(context.Background(), id)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}
