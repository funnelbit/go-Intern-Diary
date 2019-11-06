package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiaryApp_CreateNewToken(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	name := "test name " + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)

	user, err := app.FindUserByName(name)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestDiaryApp_LoginUser(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	name := "test name" + randomString()
	password := randomString() + randomString()
	err := app.CreateNewUser(name, password)
	assert.NoError(t, err)

	login, err := app.LoginUser(name, password)
	assert.NoError(t, err)
	assert.True(t, login)

	login, err = app.LoginUser(name, password+".")
	assert.NoError(t, err)
	assert.False(t, login)
}
