package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// Server : master
type Server struct {
	Router *gin.Engine
}

func (app *Server) reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("request URL: " + c.Request.URL.String())
		fmt.Println("target: " + target)

		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   target,
		})

		proxy.Director = func(r *http.Request) {
			r.Host = target
			r.URL.Host = r.Host
			r.URL.Scheme = "http"
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (app *Server) addCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Accept-Encoding, Accept-Language, Referer, User-Agent, Connection, Host")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		c.Next()
	}
}

func (app *Server) initializeRouter() {
	app.Router = gin.Default()

	app.Router.Use(app.addCORSMiddleware())

	backendURL := os.Getenv("BACK_SERVICE_HOST") + ":" + os.Getenv("BACK_SERVICE_PORT")

	// Serves the static application
	app.Router.Static("/static", "front/static")
	app.Router.StaticFile("/", "front/index.html")

	graphql := app.Router.Group("/graphql")
	{
		// Redirect /graphql/* to the backend
		graphql.GET("/*anything", app.reverseProxy(backendURL))
	}

	subscriptions := app.Router.Group("/subscriptions")
	{
		// http://localhost/subscriptions/* => http://{backendURL}/subscriptions/*
		subscriptions.GET("/*anything", app.reverseProxy(backendURL))
	}
}

// Initialize : initializes the app
func (app *Server) Initialize() {
	app.initializeRouter()
}

// Run : runs the applications
func (app *Server) Run(addr string) {
	app.Router.Run(addr)
}
