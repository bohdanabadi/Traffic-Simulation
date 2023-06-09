package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	engine *gin.Engine
	config Config
}

func NewServer(cfg Config) *Server {
	srv := &Server{
		engine: gin.Default(),
		config: cfg,
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	srv.engine.Use(cors.New(corsConfig))
	return srv
}

func (srv *Server) StartServer() error {
	var err error
	api := srv.engine.Group("/api")
	{
		api.GET("/fe", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"message": "Hi!, Component data not found. Time to useState!"})
		})
	}

	// Here you could add a switch to start the server with TLS or without depending on the configuration
	env := os.Getenv("ENV")
	if env == "dev" || env == "" {
		env = "development"
	}
	switch env {
	case "development":
		err = srv.engine.Run(srv.config.ServerDev.Host + ":" + srv.config.ServerDev.Port)
	case "production":
		srv.engine.Static("/static", "./build/static")
		// Catch-all route to serve React App
		srv.engine.NoRoute(func(c *gin.Context) {
			// check if the request path starts with /api, if so, return a 404
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
				return
			}
			// otherwise, serve the index.html file
			c.File("./build/index.html")
		})
		err = srv.engine.RunTLS(srv.config.ServerProd.Host+":"+srv.config.ServerProd.Port, srv.config.CertFile, srv.config.KeyFile)
	default:
		log.Fatalf("Invalid environment: %s", env)
	}
	return err
}
