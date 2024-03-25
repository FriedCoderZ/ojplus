package controllers

import (
	"Alarm/internal/web/forms"
	"Alarm/internal/web/services"
	"log"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	svc *services.Auth
	cfg map[string]interface{}
}

func NewAuthController(cfg map[string]interface{}) *Auth {
	svc := services.NewAuth(cfg)
	return &Auth{svc: svc, cfg: cfg}
}

func (ctrl *Auth) LoginMiddleware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		response(ctx, 203, nil)
		ctx.Abort()
		return
	}

	claims, err := ctrl.svc.ValidateToken(token)
	if err != nil || claims == nil {
		response(ctx, 203, nil)
		ctx.Abort()
		return
	}
	ctx.Set("claims", claims)
	ctx.Next()
}

func (ctrl *Auth) Test(ctx *gin.Context) {
	response(ctx, 200, ctx.Value("claims"))
}

func (ctrl *Auth) Login(ctx *gin.Context) {
	form, err := forms.NewLogin(ctx)
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
	userID, err := ctrl.svc.VerifyPassword(form.Username, form.Password)
	if err != nil {
		response(ctx, 50001, nil)
		return
	}
	if userID == 0 {
		response(ctx, 203, nil)
		return
	}
	form.Model.ID = userID
	token, err := ctrl.svc.RefreshToken(form.Model, "", ctrl.cfg["tokenValidSeconds"].(int))
	if err != nil {
		log.Println(err)
		response(ctx, 50001, nil)
		return
	}
	data := map[string]interface{}{
		"token":  token,
		"userId": userID,
	}
	response(ctx, 200, data)
}
func (ctrl *Auth) Logout(ctx *gin.Context) {

}
