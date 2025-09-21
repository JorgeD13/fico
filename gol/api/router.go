package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"fico/gol/controller"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(gdb *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"https://127.0.0.1:9090", "https://127.0.0.1:3000",
		"https://localhost:3000", "https://localhost:8080", "https://localhost:9090",
		"http://127.0.0.1:3000", "http://localhost:3000",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Recovery en español
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("API Panic Recovered: %v\nStack Trace: %s", recovered, debug.Stack())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lo sentimos, ha ocurrido un error inesperado en el servidor",
		})
	}))

	// Health
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"detail": "success"}) })

	// Auth public
	r.POST("/login", gin.WrapF(controller.LoginHandler(gdb)))
	r.POST("/logout", gin.WrapF(controller.Logout(gdb)))

	// Protected
	auth := r.Group("/")
	auth.Use(JWTMiddleware(gdb))
	auth.GET("/getUser", gin.WrapF(controller.GetUser(gdb)))
	auth.GET("/getUsers", gin.WrapF(controller.GetUsers(gdb)))
	auth.POST("/createUser", gin.WrapF(controller.CreateUser(gdb)))
	auth.POST("/editUser", gin.WrapF(controller.EditUser(gdb)))
	auth.POST("/deleteUser", gin.WrapF(controller.DeleteUser(gdb)))

	return r
}

func SetUpServer(gdb *gorm.DB) {
	fmt.Printf("Starting server...")
	r := SetUpRouter(gdb)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
