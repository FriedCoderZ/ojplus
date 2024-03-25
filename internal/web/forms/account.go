package forms

import (
	"Alarm/internal/web/models"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserCreate struct {
	Username string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=6,max=128"`
	Name     string `validate:"required,max=24"`
	Email    string `validate:"omitempty,email"`
	Phone    string `validate:"omitempty,number,min=6,max=32"`
	Model    *models.User
}

func NewUserCreate(ctx *gin.Context) (*UserCreate, error) {
	var form *UserCreate
	err := ctx.BindJSON(&form)
	if err != nil {
		return nil, err
	}
	form.Model = &models.User{
		Username: form.Username,
		Password: form.Password,
		Name:     form.Name,
		Email:    form.Email,
		Phone:    form.Phone,
	}
	return form, nil
}

func (form *UserCreate) check() map[string]string {
	result := make(map[string]string)
	if form.Phone != "" && !checkPhone(form.Phone) {
		result["phone"] = "电话号码格式错误"
	}
	return result
}

func checkPhone(phone string) bool {
	regex := `^(?:\+)?[0-9]{1,3}[-.●]?\(?[0-9]{1,3}\)?[-.●]?[0-9]{1,4}[-.●]?[0-9]{1,4}$`
	pattern := regexp.MustCompile(regex)
	matched := pattern.MatchString(phone)
	return matched
}
