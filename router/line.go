package router

import (
	"context"
	"fmt"
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
	"github.com/spf13/viper"
)

var lineClient *linebot.Client
var domain, intr, aboutUs string

func NewLineRouter(group *gin.RouterGroup, bot *linebot.Client, cnf *viper.Viper) {
	lineClient = bot
	domain = cnf.GetString("frontend.host")
	intr = cnf.GetString("line.intr")
	aboutUs = cnf.GetString("line.aboutUs")
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

				switch input {
				case "操作說明":
					jsonData := []byte(intr)

					container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
					// err is returned if invalid JSON is given that cannot be unmarshalled
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
						return
					}
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("操作說明", container)).Do(); err != nil {
						return
					}

				case "當前設定":
					userSetting, err := user.UserService.GetLineUserSetting(ctx, event.Source.UserID)
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
						return
					}

					result := fmt.Sprintf("媒體檔案可瀏覽秒數: %d\n媒體檔案瀏覽密碼: %s", userSetting.ExpirationTime, userSetting.Password)
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						return
					}

				case "關於我們":
					jsonData := []byte(aboutUs)

					container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
					// err is returned if invalid JSON is given that cannot be unmarshalled
					if err != nil {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
							return
						}
						return
					}
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("關於我們", container)).Do(); err != nil {
						return
					}

				default:

					index := strings.Index(input, ":")
					if index == -1 {
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
							return
						}
						return
					}

					switch input[:index] {
					case "set time":
						input = input[index+1:]
						expirationTime, err := strconv.ParseInt(input, 10, 64)
						if err != nil {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
								return
							}
							return
						}

						if expirationTime <= 0 || expirationTime > 604800 {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
								return
							}
							return
						}

						userSetting, err := user.UserService.UpdateLineUserExpirationTime(ctx, event.Source.UserID, expirationTime)
						if err != nil {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
								return
							}
							return
						}

						result := fmt.Sprintf("成功設定媒體檔案可瀏覽秒數: %d", userSetting.ExpirationTime)

						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
							return
						}

					case "set pass":
						input = input[index+1:]
						if input != "none" && input != "today" {
							_, err := strconv.ParseInt(input, 10, 64)
							if err != nil {
								if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
									return
								}
								return
							}
						}

						if input == "today" {
							today := time.Now().Format("2006-01-02")
							t := strings.Split(today, "-")
							input = t[1] + t[2]
						}

						if len(input) > 10 {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
								return
							}
							return
						}

						userSetting, err := user.UserService.UpdateLineUserPassword(ctx, event.Source.UserID, input)
						if err != nil {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
								return
							}
							return
						}

						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功設定媒體檔案瀏覽密碼: "+userSetting.Password)).Do(); err != nil {
							return
						}

					case "https", "http":
						info := ShortInfo{}
						info.LeadUrl = message.Text

						validate := validator.New()
						err = validate.Struct(info)
						if err != nil {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
								return
							}
							return
						}

						shortUrl := hash.StringHash(info.LeadUrl)

						data, err := short.ShortService.CreateShort(ctx, event.Source.UserID, shortUrl, info.LeadUrl)
						if err != nil {
							if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
								return
							}
							return
						}

						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(domain+data.ShortUrl)).Do(); err != nil {
							return
						}

					default:
						if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無效的輸入(๑╹◡╹๑)")).Do(); err != nil {
							return
						}
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
					return
				}

				_, buf, err := fileHelper.DownscaleImageDefault(byte)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
					return
				}

				userSetting, err := user.UserService.GetLineUserSetting(ctx, event.Source.UserID)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
				}

				data, err := media.MediaService.CreateMedia(ctx, event.Source.UserID, "image", userSetting.Password, userSetting.ExpirationTime, buf.Bytes())
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
					return
				}

				if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(domain+data.ShortUrl)).Do(); err != nil {
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
					return
				}

				userSetting, err := user.UserService.GetLineUserSetting(ctx, event.Source.UserID)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
					return
				}

				data, err := media.MediaService.CreateMedia(ctx, event.Source.UserID, "video", userSetting.Password, userSetting.ExpirationTime, byte)
				if err != nil {
					if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("發生未知錯誤∑(✘Д✘๑ )")).Do(); err != nil {
						return
					}
					return
				}

				if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(domain+data.ShortUrl)).Do(); err != nil {
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
