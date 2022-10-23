package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/controller"
)

// Route makes the routing
func Route(app *gin.Engine) {
	indexController := new(controller.IndexController)
	searchTwitterController := new(controller.SearchTwitterController)

	app.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error400.html", gin.H{
			"title":   "Error",
			"content": "No encontrada",
		})
	})

	api := app.Group("/api")
	api.GET("/version", indexController.GetVersion)
	api.GET("/search/twitter", searchTwitterController.Search)
}
