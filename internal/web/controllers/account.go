package controllers

import (
	"Alarm/internal/web/forms"
	"Alarm/internal/web/models"
	"Alarm/internal/web/services"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Account struct {
	svc *services.Account
	cfg map[string]interface{}
}

func NewAccount(cfg map[string]interface{}) *Account {
	svc := services.NewAccount(cfg)
	return &Account{svc: svc, cfg: cfg}
}

func (ctrl *Account) CreateUser(ctx *gin.Context) {
	form, err := forms.NewUserCreate(ctx)
	if err != nil {
		response(ctx, 40001, nil)
		return
	}
	isValid, errorsMap, err := forms.Verify(form)
	if err != nil {
		response(ctx, 50001, nil)
		return
	}
	if !isValid {
		response(ctx, 400, errorsMap)
		return
	}
	user := form.Model
	has, hasMessage, err := ctrl.svc.IsUserExist(user)
	if err != nil {
		response(ctx, 50001, nil)
		return
	}
	if has {
		responseWithMessage(ctx, hasMessage, 40901, nil)
		return
	}
	err = ctrl.svc.CreateUser(user)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1062 {
			response(ctx, 40901, nil)
			return
		}
	} else if err != nil {
		response(ctx, 400, nil)
	}
	response(ctx, 201, map[string]int{"userId": user.ID})
}

func (ctrl *Account) AllUser(ctx *gin.Context) {
	users, err := ctrl.svc.AllUser()
	if err != nil {
		response(ctx, 500, nil)
		return
	}
	response(ctx, 200, users)
}

func (ctrl *Account) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response(ctx, 400, nil)
		return
	}
	var userInfo *models.UserInfo
	userInfo, err = ctrl.svc.GetUserByID(idInt)
	if err != nil {
		response(ctx, 500, nil)
		log.Println(err)
		return
	}
	if userInfo == nil || userInfo.ID == 0 {
		response(ctx, 404, nil)
		return
	}
	response(ctx, 200, userInfo)
}

func (ctrl *Account) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response(ctx, 400, nil)
		return
	}
	var user *models.User
	err = ctx.BindJSON(&user)
	if err != nil {
		response(ctx, 40001, nil)
		return
	}
	user.ID = 0
	ctrl.svc.UpdateUserByID(idInt, user)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1062 {
			response(ctx, 40901, nil)
			return
		}
	}
}
