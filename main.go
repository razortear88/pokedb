package main

import (
	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/routes" //add this
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/assets", "./assets")
	//run database
	configs.ConnectDB()

	//routes
	routes.MainRoute(router) //add this

	router.Run("localhost:8080")
}
