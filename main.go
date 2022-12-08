package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maickmachado/golang-ecommerce/controllers"
	"github.com/maickmachado/golang-ecommerce/database"
	"github.com/maickmachado/golang-ecommerce/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use()

	log.Fatal(router.Run(":" + port))

}
