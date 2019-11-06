package service

import (
	"testing"
	"time"

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

	token, err := app.CreateNewToken(user.ID, time.Now().Add(1*time.Hour))
	assert.NoError(t, err)
	assert.NotEqual(t, "", token)

	u, err := app.FindUserByToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, u.ID)
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
