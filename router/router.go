package router

import (
	"time"

	"github.com/Elvis-Benites-N/GolangChat/internal/user"
	"github.com/Elvis-Benites-N/GolangChat/internal/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	r = gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Let's start": "Great! this Chat service is working - use Postman to create an account in  http://localhost:4200/signup",
		})
	})
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)
	r.POST("/ws/createServer", wsHandler.CreateServer)
	r.GET("/ws/joinServer/:serverId", wsHandler.JoinServer)
	r.GET("/ws/getServers", wsHandler.GetServers)
	r.GET("/ws/getClients/:serverId", wsHandler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
