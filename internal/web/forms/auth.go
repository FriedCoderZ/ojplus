package forms

import (
	"Alarm/internal/web/models"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Model    *models.User
}

func NewLogin(ctx *gin.Context) (*Login, error) {
	var form *Login
	err := ctx.BindJSON(&form)
	if err != nil {
		return nil, err
	}
	form.Model = &models.User{
		Username: form.Username,
		Password: form.Password,
	}
	return form, nil
}

func (form *Login) check() map[string]string {
	result := make(map[string]string)
	return result
}
