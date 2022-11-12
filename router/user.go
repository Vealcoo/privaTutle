package router

import (
	"context"
	"net/http"
	"privaTutle/model"
	"privaTutle/pkg/auth"
	"privaTutle/service/short"
	"privaTutle/service/user"
	"strconv"
	"time"

	httpHelper "privaTutle/pkg/http_helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func NewUserRouter(group *gin.RouterGroup) {
	group.POST("/register", Register)
	group.POST("/login", Login)
	group.GET("/short/:page/:limit", ShortList)
	group.DELETE("/short/:shortId", DeleteShort)
}

type RegisterInfo struct {
	UserId       string `validate:"required,min=6,max=20"`
	UserPassword string `validate:"required,min=6,max=20"`
	UserName     string `validate:"required,min=2,max=20"`
	Email        string `validate:"required,email"`
	Sex          string `validate:"required,oneof=female male none"`
}

// @Summary Register
// @Tags User
// @Accept  json
// @produce json
// @Param  body  body  RegisterInfo  true  "body"
// @Success 200
// @Router /api/user/register [post]
func Register(g *gin.Context) {

	info := RegisterInfo{}
	g.BindJSON(&info)

	validate := validator.New()
	err := validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := user.UserService.Register(ctx, info.UserId, info.UserPassword, info.UserName, info.Email, info.Sex)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	token, err := auth.SetToken(data.Id)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	httpHelper.SendResponse(g, gin.H{
		"data":  data,
		"token": token,
	})
}

type LoginInfo struct {
	UserId       string `validate:"required,min=6,max=20"`
	UserPassword string `validate:"required,min=6,max=20"`
}

// @Summary Login
// @Tags User
// @Accept  json
// @produce json
// @Param  body  body  LoginInfo  true  "body"
// @Success 200
// @Router /api/user/login [post]
func Login(g *gin.Context) {

	info := LoginInfo{}
	g.BindJSON(&info)

	validate := validator.New()
	err := validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := user.UserService.Login(ctx, info.UserId, info.UserPassword)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	token, err := auth.SetToken(data.Id)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	// TODO: set token in redis or temp collection

	httpHelper.SendResponse(g, gin.H{
		"data":  data,
		"token": token,
	})
}

type ShortListInfo struct {
	Page  int64 `validate:"required,gte=1"`
	Limit int64 `validate:"required,gte=1,lte=20"`
}

// @Summary ShortList
// @Tags User
// @Accept  json
// @produce json
// @Param  Authorization  header  string  true  "Authorization"
// @Param  page  path  int64  true  "page"
// @Param  limit  path  int64  true  "limit"
// @Success 200
// @Router /api/user/short/{page}/{limit} [get]
func ShortList(g *gin.Context) {
	token := g.Request.Header.Get("Authorization")
	objectId, err := auth.AuthJWT(token)
	if err != nil {
		if err == auth.ErrVaild {
			httpHelper.SendError(g, http.StatusUnauthorized, err.Error())
			return
		}
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	page, err := strconv.ParseInt(g.Param("page"), 10, 64)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}
	limit, err := strconv.ParseInt(g.Param("limit"), 10, 64)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}
	info := ShortListInfo{
		Page:  page,
		Limit: limit,
	}

	validate := validator.New()
	err = validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := short.ShortService.ListUserShorts(ctx, objectId, info.Page, info.Limit)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, data)
}

type DeleteShortInfo struct {
	ShortId string `validate:"required"`
}

// @Summary DeleteShort
// @Tags User
// @Accept  json
// @produce json
// @Param  Authorization  header  string  true  "Authorization"
// @Param  shortId  path  string  true  "shortId"
// @Success 200
// @Router /api/user/short/{shortId} [delete]
func DeleteShort(g *gin.Context) {
	token := g.Request.Header.Get("Authorization")
	objectId, err := auth.AuthJWT(token)
	if err != nil {
		if err == auth.ErrVaild {
			httpHelper.SendError(g, http.StatusUnauthorized, err.Error())
			return
		}
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	info := DeleteShortInfo{
		ShortId: g.Param("shortId"),
	}

	validate := validator.New()
	err = validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = short.ShortService.UpdateShortStatus(ctx, objectId, info.ShortId, "delete")
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, nil)
}
