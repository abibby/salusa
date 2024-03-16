package factories

import (
	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/go-faker/faker/v4"
)

var UserFactory = dbtest.NewFactory(func() *models.User {
	return &models.User{
		BaseUser: auth.BaseUser{
			Username:     faker.Username(),
			PasswordHash: []byte(faker.Password()),
		},
	}
})
