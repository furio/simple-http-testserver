package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Delay  int
	Mirror bool
	Cors   bool
}

// GenerateHTTPRoutes is a function that creates a Gin Engine according to the configuration
func GenerateHTTPRoutes(config RouterConfig) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if config.Cors {
		r.Use(cors.Default())
	}

	// Routes
	r.Any("/*any", func(c *gin.Context) {
		bodyData := "Hello from test server"
		if config.Mirror {
			byteData, _ := c.GetRawData()
			bodyData = string(byteData)
		}

		if config.Delay > 0 {
			time.Sleep(time.Duration(config.Delay) * time.Millisecond)
		}

		c.JSON(200, gin.H{
			"message": bodyData,
		})
	})

	return r
}
