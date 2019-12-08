package dao

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/internal/test"
)

func TestMain(m *testing.M) {
	logger.Init(test.GetLogLevel())

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestNewUser(t *testing.T) {
	dbHandler, err := NewHandler(
		"sqlite3", ":memory:", []CompactSample{})
	assert.NoError(t, err, "DB connect")
	defer dbHandler.Close() //nolint:wsl

	samples := GetDefaultSampleFill()
	user0 := samples[0].User

	user, err := dbHandler.CreateUser(user0.Name, user0.Password)
	assert.NoError(t, err, "Insert User")
	assert.Equal(t, user0.Name, user.Name, "User name")
	assert.Equal(t, hashPassword(user0.Password), user.Password, "User password")
}

func TestGetUser(t *testing.T) {
	samples := GetDefaultSampleFill()
	user0 := samples[1]

	dbHandler, err := NewHandler(
		"sqlite3", ":memory:", GetDefaultSampleFill())
	assert.NoError(t, err, "DB connect")
	defer dbHandler.Close() //nolint:wsl

	user, err := dbHandler.GetUserByName(user0.Name)
	assert.NoError(t, err, "Get User")
	assert.Equal(t, user0.Name, user.Name, "User name")
	assert.Equal(t, hashPassword(user0.Password), user.Password, "User password")
}

func TestGetUserFailed(t *testing.T) {
	samples := GetDefaultSampleFill()

	dbHandler, err := NewHandler(
		"sqlite3", ":memory:", samples)
	assert.NoError(t, err, "DB connect")
	defer dbHandler.Close() //nolint:wsl

	_, err = dbHandler.GetUserByName("Non noN")
	assert.Error(t, err, "Get User")
}
