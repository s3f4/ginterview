package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InMemoryRepository(t *testing.T) {
	inMemoryRepository := NewInMemoryRepository()
	inMemoryRepository.Create("active-tabs", "getir")
	assert.Equal(t, inMemoryRepository.Exist("active-tabs"), true)
	assert.Equal(t, inMemoryRepository.Exist("active-tab"), false)
	assert.Equal(t, inMemoryRepository.Get("active-tabs"), "getir")
	assert.Nil(t, inMemoryRepository.Get("active-tab"))
}
