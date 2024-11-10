package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainRoute(app *gin.Engine) {
	route := app.Group("")

	TypeRoute(route)
	AbilityRoute(route)
	GameRoute(route)
	MoveRoute(route)
	route.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Create Typesss",
		})
	})
	// AbilityRoute(route)
}
