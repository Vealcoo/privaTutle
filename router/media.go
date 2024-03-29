package router

import (
	"bytes"
	"context"
	"net/http"
	"privaTutle/model"
	"privaTutle/pkg/auth"
	"privaTutle/service/media"
	"strconv"
	"time"

	fileHelper "privaTutle/pkg/file_helper"
	httpHelper "privaTutle/pkg/http_helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func NewMediaRouter(group *gin.RouterGroup) {
	group.POST("/image", UploadImage)
	group.POST("/video", UploadVideo)
	group.GET("/:short", GetMedia)
}

type UploadMediaInfo struct {
	ExpirationTime int64  `validate:"required,gte=1,lte=86400"`
	Password       string `validate:"max=10"`
}

// @Summary UploadImage
// @Tags Media
// @Accept  mpfd
// @produce json
// @Param  Authorization  header  string  false  "Authorization"
// @Param  image  formData  file  true  "上傳圖片"
// @Param  expirationTime  formData  string  true  "有效時間"
// @Param  password  formData  string  false  "瀏覽密碼"
// @Success 200
// @Router /api/media/image [post]
func UploadImage(g *gin.Context) {
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

	file, err := g.FormFile("image")
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, err.Error())
		return
	}

	var buf *bytes.Buffer
	if file != nil {
		b, err := fileHelper.ReadFile(file)
		if err != nil {
			httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
			return
		}
		if len(b) == 0 {
			return
		}
		if !fileHelper.IsImage(http.DetectContentType(b)) {
			httpHelper.SendError(g, http.StatusBadRequest, "ErrInvalidFileType")
			return
		}

		_, buf, err = fileHelper.DownscaleImageDefault(b)
		if err != nil {
			httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
			return
		}
	} else {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	expirationTime, err := strconv.ParseInt(g.PostForm("expirationTime"), 10, 64)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	info := &UploadMediaInfo{
		ExpirationTime: expirationTime,
		Password:       g.PostForm("password"),
	}

	validate := validator.New()
	err = validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	password := info.Password
	if password == "" {
		password = "none"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := media.MediaService.CreateMedia(ctx, objectId, "image", password, info.ExpirationTime, buf.Bytes())
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, data)
}

// @Summary UploadVideo
// @Tags Media
// @Accept  mpfd
// @produce json
// @Param  Authorization  header  string  false  "Authorization"
// @Param  video  formData  file  true  "上傳影片"
// @Param  expirationTime  formData  string  true  "有效時間"
// @Param  password  formData  string  false  "瀏覽密碼"
// @Success 200
// @Router /api/media/video [post]
func UploadVideo(g *gin.Context) {
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

	file, err := g.FormFile("video")
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, err.Error())
		return
	}

	var buf []byte
	if file != nil {
		buf, err = fileHelper.ReadFile(file)
		if err != nil {
			httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
			return
		}
		if len(buf) == 0 {
			return
		}
		if !fileHelper.IsVideo(http.DetectContentType(buf)) {
			httpHelper.SendError(g, http.StatusBadRequest, "ErrInvalidFileType")
			return
		}
	} else {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	expirationTime, err := strconv.ParseInt(g.PostForm("expirationTime"), 10, 64)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	info := &UploadMediaInfo{
		ExpirationTime: expirationTime,
		Password:       g.PostForm("password"),
	}

	validate := validator.New()
	err = validate.Struct(info)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, model.ErrParameter.Error())
		return
	}

	password := info.Password
	if password == "" {
		password = "none"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := media.MediaService.CreateMedia(ctx, objectId, "video", password, info.ExpirationTime, buf)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, data)
}

// @Summary GetMedia
// @Tags Media
// @Accept  json
// @produce json
// @Param  Authorization  header  string  false  "Authorization"
// @Param  short  path  string  true  "short"
// @Param  password  query  string  false  "password"
// @Success 200
// @Router /api/media/{short} [get]
func GetMedia(g *gin.Context) {
	shortUrl := g.Param("short")
	password := g.Query("password")

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := media.MediaService.TranslateMedia(ctx, shortUrl, password, objectId)
	if err != nil {
		httpHelper.SendError(g, http.StatusBadRequest, err.Error())
		return
	}

	httpHelper.SendResponse(g, data)
}
