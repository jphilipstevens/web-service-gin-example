package main

// @title Web Service Gin Example
// @version 1.0
// @description Example API demonstrating middleware and Swagger docs.
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	"flag"
	"os"

	"github.com/jphilipstevens/web-service-gin-example/example/features/albums"
	"github.com/jphilipstevens/web-service-gin-example/example/seed"
	server "github.com/jphilipstevens/web-service-gin/v2/pkg"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/cache"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/config"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/db"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/dependencies"

	"github.com/jphilipstevens/web-service-gin-example/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(deps *dependencies.Dependencies[db.Database]) {
	albums.Init(deps)
	deps.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RunApp() {
	docs.SwaggerInfo.BasePath = "/"

	srv, err := server.NewServer[db.Database](config.ConfigOptions{
		Path: "./example/config", // or from ENV, flags, etc
		Name: "config",           // without extension
		Type: "yaml",             // optional
	})
	if err != nil {
		panic(err)
	}

	cfg := srv.Config()
	tracer := srv.Dependencies().Tracer

	dbConn, err := db.NewDatabase(cfg.DB, tracer)
	if err != nil {
		panic(err)
	}
	cacheClient := cache.NewCacher(cfg.Redis, tracer)

	srv.WithDatabase(dbConn)
	srv.WithCache(cacheClient)

	srv.RegisterRoutes(registerRoutes)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "seed":
			seed.Init()
			os.Exit(0)
		default:
			RunApp()
		}
	} else {
		RunApp()
	}
}
