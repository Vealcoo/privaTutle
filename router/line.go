package router

import (
	"context"
	"io/ioutil"
	"privaTutle/model"
	fileHelper "privaTutle/pkg/file_helper"
	"privaTutle/pkg/hash"
	httpHelper "privaTutle/pkg/http_helper"
	"privaTutle/service/media"
	"privaTutle/service/short"
	"privaTutle/service/user"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var lineClient *linebot.Client

func NewLineRouter(group *gin.RouterGroup, bot *linebot.Client) {
	lineClient = bot
	group.POST("", LineCallback)
}

func LineCallback(g *gin.Context) {
	events, err := lineClient.ParseRequest(g.Request)
	if err != nil {
		httpHelper.SendError(g, http.StatusInternalServerError, model.ErrInternal.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				input := message.Text
				index := strings.Index(input, "/")
				switch input[:index] {

				case "setTime":
					input = input[index+1:]
					expirationTime, err := strconv.ParseInt(input, 10, 64)
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
						return
					}

					if expirationTime <= 0 || expirationTime >= 86400 {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
							return
						}
					}

					userSetting, err := user.UserService.UpdateLineUserExpirationTime(ctx, event.Source.UserID, expirationTime)
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
					}

					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功設定媒體檔案可瀏覽秒數: "+strconv.FormatInt(userSetting.ExpirationTime, 10))).Do(); err != nil {
					}

				case "setPass":
					input = input[index+1:]
					_, err := strconv.ParseInt(input, 10, 64)
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
						return
					}

					if len(input) < 4 || len(input) > 10 {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
							return
						}
					}

					userSetting, err := user.UserService.UpdateLineUserPassword(ctx, event.Source.UserID, input)
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
					}

					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功設定媒體檔案瀏覽密碼: "+userSetting.Password)).Do(); err != nil {
					}

				case "getSet":
					// TODO: 取得使用者當前設定資料
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("")).Do(); err != nil {
					}

				case "https:", "http:":
					info := ShortInfo{}
					info.LeadUrl = message.Text

					validate := validator.New()
					err = validate.Struct(info)
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
							return
						}
					}

					shortUrl := hash.StringHash(info.LeadUrl)

					data, err := short.ShortService.CreateShort(ctx, event.Source.UserID, strconv.FormatInt(int64(shortUrl), 10), info.LeadUrl)
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
					}

					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("網域/"+data.ShortUrl)).Do(); err != nil {
						return
					}

				default:
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
						return
					}
				}

			case *linebot.ImageMessage:
				content, err := lineClient.GetMessageContent(message.ID).Do()
				if err != nil {

				}
				defer content.Content.Close()

				byte, err := ioutil.ReadAll(content.Content)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				_, buf, err := fileHelper.DownscaleImageDefault(byte)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				userSetting, err := user.UserService.GetLineUserSetting(ctx, event.Source.UserID)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				data, err := media.MediaService.CreateMedia(ctx, event.Source.UserID, "image", "", userSetting.ExpirationTime, buf.Bytes())
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("網域/"+data.ShortUrl)).Do(); err != nil {
					return
				}

			case *linebot.VideoMessage:

				content, err := lineClient.GetMessageContent(message.ID).Do()
				if err != nil {

				}
				defer content.Content.Close()

				byte, err := ioutil.ReadAll(content.Content)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				userSetting, err := user.UserService.GetLineUserSetting(ctx, event.Source.UserID)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				data, err := media.MediaService.CreateMedia(ctx, event.Source.UserID, "video", "", userSetting.ExpirationTime, byte)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("網域/"+data.ShortUrl)).Do(); err != nil {
					return
				}

			default:
				if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
					return
				}

			}
		}
	}
}
