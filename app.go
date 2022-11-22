package privatutle

import (
	"context"
	"os"
	"privaTutle/router"

	"privaTutle/service/media"
	"privaTutle/service/short"
	"privaTutle/service/user"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cnf *viper.Viper

func configInit() {
	cnf = viper.New()
	cnf.AddConfigPath("./config")
	cnf.SetConfigName("app")
	cnf.SetConfigType("yaml")
	cnf.AutomaticEnv()

	err := cnf.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func logInit() {

}

func dbConn() *mongo.Database {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cnf.GetString("mongo.applyURI")))
	if err != nil {
		panic(err)
	}

	return mongoClient.Database("privaTutle")
}

func gcsConn() *storage.Client {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cnf.GetString("google.credentials"))
	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return gcsClient
}

func lineBotConn() *linebot.Client {
	bot, err := linebot.New(cnf.GetString("line.secret"), cnf.GetString("line.accessToken"))
	if err != nil {
		panic(err)
	}

	return bot
}

func serviceBuild(database *mongo.Database, gcsClient *storage.Client) {
	user.NewUserService(database)
	short.NewShortService(database)
	media.NewMediaService(database, gcsClient)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		g.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		g.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		g.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if g.Request.Method == "OPTIONS" {
			g.AbortWithStatus(204)
			return
		}

		g.Next()
	}
}

func Run() {
	configInit()
	logInit()
	database := dbConn()
	gcsClient := gcsConn()
	botClient := lineBotConn()
	serviceBuild(database, gcsClient)

	g := gin.Default()
	g.Use(CORSMiddleware())
	router.NewUserRouter(g.Group("api/user"))
	router.NewMediaRouter(g.Group("api/media"))
	router.NewShortRouter(g.Group("api/short"))
	router.NewLineRouter(g.Group("api/line"), botClient, cnf)
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g.Run(":8888")
}
