package main

import (
	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/routes" //add this
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	//run database
	configs.ConnectDB()

	//routes
	routes.TypeRoute(router) //add this

	router.Run("localhost:8080")
}
