package router

import (
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
		if config.Mirror {
			bodyData, _ := c.GetRawData()

			c.JSON(200, gin.H{
				"message": string(bodyData),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Hello from test server",
			})
		}
	})

	return r
}
