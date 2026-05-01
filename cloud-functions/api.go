package main

import (
	"bbs-go/internal/install"
	"bbs-go/internal/router"
	"bbs-go/internal/services/eventhandler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	install.InitConfig()
	install.InitLogger()
	install.InitLocales()

	if install.IsInstalled() {
		if err := install.InitDB(); err != nil {
			panic(err)
		}
		if err := install.InitMigrations(); err != nil {
			panic(err)
		}
		install.InitOthers()
	}

	_ = eventhandler.RegisterHandlers

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	router.Setup(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
