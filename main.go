package main

import (
	"fmt"
	"log"

	"github.com/arthur-rc18/Go-Redis/database/scripts"

	"github.com/arthur-rc18/Go-Redis/routes"

	"github.com/gin-gonic/gin"
)

const port string = "localhost:8000"

func main() {

	// redis := database.ConnectRedis()
	// defer redis.Close()
	scripts.UpdateDatabase()
	scripts.PopulateDatabase(nil)

	router := gin.Default()
	routes.StartRoutes(router)

	if err := router.Run(port); err != nil {
		err := fmt.Errorf("Could not run the application: %v", err)
		log.Fatalf(err.Error())
	}
}
