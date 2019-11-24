package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/controller"
	"github.com/hyperjiang/gin-skeleton/middleware"
	"github.com/hyperjiang/gin-skeleton/model"
)

// Route makes the routing
func Route(app *gin.Engine) {
	indexController := new(controller.IndexController)
	userController := new(controller.UserController)
	authMiddleware := middleware.Auth()

	app.GET(
		"/", indexController.GetIndex,
	).GET(
		"/user/:id", userController.GetUser,
	).GET(
		"/signup", userController.SignupForm,
	).POST(
		"/signup", userController.Signup,
	).GET(
		"/login", userController.LoginForm,
	).POST(
		"/login", authMiddleware.LoginHandler,
	)

	app.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error400.html", gin.H{
			"title":   "Error",
			"content": "No encontrada",
		})
	})

	auth := app.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get("email")
			c.JSON(200, gin.H{
				"email": claims["email"],
				"name":  user.(*model.User).Name,
				"text":  "Hello World.",
			})
		})
	}

	api := app.Group("/api")
	api.GET("/version", indexController.GetVersion)
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get("email")
			c.JSON(200, gin.H{
				"email": claims["email"],
				"name":  user.(*model.User).Name,
				"text":  "Hello World.",
			})
		})
	}
}
