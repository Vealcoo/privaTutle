package router

import (
	"net/http"
	"privaTutle/model"
	"privaTutle/pkg/auth"
	"privaTutle/pkg/hash"
	httpHelper "privaTutle/pkg/http_helper"
	"privaTutle/service/short"
	"strconv"

	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func NewShortRouter(group *gin.RouterGroup) {
	group.POST("", Short)
	group.GET("/:short", GetShort)
}

type ShortInfo struct {
	LeadUrl string `validate:"required,url"`
}

// @Summary Short
// @Tags Short
// @Accept  json
// @produce json
// @Param  Authorization  header  string  false  "Authorization"
// @Param  body  body  ShortInfo  true  "body"
// @Success 200
// @Router /api/short [post]
func Short(g *gin.Context) {
	token := g.Request.Header.Get("Authorization")
	var objectId string
	var err error
	if token != "" {
		objectId, err = auth.AuthJWT(token)
		if err != nil {
			if err == auth.ErrVaild {
				httpHelper.SendError(g, http.StatusUnauthorized, err.Error())
				return
			}
			httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
			return
		}
	}

	info := ShortInfo{}
	g.BindJSON(&info)

	validate := validator.New()
	err = validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	shortUrl := hash.StringHash(info.LeadUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := short.ShortService.CreateShort(ctx, objectId, strconv.FormatInt(int64(shortUrl), 10), info.LeadUrl)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, gin.H{
		"shortUrl": data.ShortUrl,
	})
}

// @Summary GetShort
// @Tags Short
// @Accept  json
// @produce json
// @Param  short  path  string  true  "short"
// @Success 200
// @Router /api/short/{short} [get]
func GetShort(g *gin.Context) {
	shortUrl := g.Param("short")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := short.ShortService.TranslateShort(ctx, shortUrl)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, gin.H{
		"leadUrl": data.LeadUrl,
	})
}
