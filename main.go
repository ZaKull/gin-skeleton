package main

import (
	"flag"
	"time"

	//"log"
	"path/filepath"

	"github.com/dvwright/xss-mw"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	_ "github.com/golang/glog"
	"github.com/hyperjiang/gin-skeleton/config"
	"github.com/hyperjiang/gin-skeleton/router"
	//"github.com/gin-gonic/autotls"
	//"golang.org/x/crypto/acme/autocert"
)

func main() {

	addr := flag.String("addr", config.Server.Addr, "Address to listen and serve")
	flag.Parse()

	if config.Server.Mode == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Secure
	app.Use(secure.New(secure.Config{
		AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
		SSLHost:               "ssl.example.com",
		IsDevelopment:         true, // cambiar cuando este en prod.
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))

	// Install nice.Recovery, passing the handler to call after recovery
	app.Use(nice.Recovery(func(c *gin.Context, err interface{}) {
		c.HTML(500, "error500.html", gin.H{
			"title":   "Error",
			"content": err,
		})
	}))

	//Xss Middleware menos campo password
	var xssMdlwr xss.XssMw
	app.Use(xssMdlwr.RemoveXss())

	app.Static("/images", filepath.Join(config.Server.StaticDir, "img"))
	app.StaticFile("/favicon.ico", filepath.Join(config.Server.StaticDir, "img/favicon.ico"))
	app.LoadHTMLGlob(config.Server.ViewDir + "/*")
	app.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20

	router.Route(app)
	/* Auto tsl, activar cuando este en ip publica.

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("zakull.test", "app.zakull.test"),
		Cache:      autocert.DirCache(config.Server.Cache),
	}

	log.Fatal(autotls.RunWithManager(app, &m))
	*/
	// Listen and Serve
	app.Run(*addr)

}
