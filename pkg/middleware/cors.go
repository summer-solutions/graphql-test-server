package middleware

import (
	"time"

	"github.com/summer-solutions/spring/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(engine *gin.Engine) error {
	configService, hasConfig := service.Config()
	if !hasConfig {
		return nil
	}

	origins := configService.GetStringSlice("cors")
	if len(origins) == 0 {
		return nil
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOrigins:     origins,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"X-Invalid-Authorization"},
	}
	engine.Use(cors.New(corsConfig))
	return nil
}
