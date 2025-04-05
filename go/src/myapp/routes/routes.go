package routes

import (
	"myapp/controller"
	"myapp/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS設定を追加
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://cooksensing.vercel.app", // 本番
			"https://*.vercel.app",           // Preview 環境用
		},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ミドルウェアを追加
	router.Use(middleware.Logger())

	// APIルートの設定
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("", controller.GetUserAll)
		userRoutes.POST("", controller.CreateUser)
		userRoutes.POST("/search_user_by_email", controller.SearchUserByEmailHandler)
		userRoutes.POST("/search_user_by_firebase_auth_uid", controller.SearchUserByFirebaseAuthUidHandler)
		userRoutes.POST("/search_user_by_name_and_firebase_auth_uid", controller.SearchUserByNameAndUIDHandler)
	}
	featureDataRoutes := router.Group("/feature_data")
	{
		featureDataRoutes.POST("", controller.CreateFeatureDatas)

		//折れ線グラフ用
		featureDataRoutes.POST("/by_userid", controller.GetFeatureDatasByUserID)

		//折れ線グラフ用（日付指定込み）daysに日数を入れる 1ヶ月なら30 1年前なら365 全部なら0 をリクエストに入れる
		featureDataRoutes.POST("/by_userid_within_days", controller.GetFeatureDatasByUserIDWithinDays)
	}
	bestRoutes := router.Group("/best")
	{
		bestRoutes.GET("", controller.GetBestAll)
		bestRoutes.GET("/average", controller.GetBestAverage)
	}

	histogramRoutes := router.Group("/histogram")
	{
		histogramRoutes.GET("", controller.GetHistogramAll)
	}

	return router
}
