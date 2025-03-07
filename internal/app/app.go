package app

import (
	"blog/config"
	"blog/internal/adapter/cloudflare"
	"blog/internal/adapter/handler"
	"blog/internal/adapter/repository"
	"blog/internal/core/service"
	"blog/lib/auth"
	"blog/lib/middleware"
	"blog/lib/pagination"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("error connecting to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatal("error creating directory: %v", err)
		return
	}

	// CloudFlare R2
	cdfR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cdfR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)

	_ = pagination.NewPagination()

	// Repository
	authRepository := repository.NewAuthRepository(db.DB)
	categoryRepository := repository.NewCategoryRepository(db.DB)
	contentRepository := repository.NewContentRepository(db.DB)
	userRepository := repository.NewUserRepository(db.DB)

	// Service
	authService := service.NewAuthService(authRepository, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepository)
	contentService := service.NewContentService(contentRepository, cfg, r2Adapter)
	userService := service.NewUserService(userRepository)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	}

	api := app.Group("/api")

	api.Post("/login", authHandler.Login)

	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	// Category route
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/create", categoryHandler.CreateCategory)
	categoryApp.Get("/:categoryId", categoryHandler.GetCategoryByID)
	categoryApp.Put("/:categoryId", categoryHandler.EditCategoryByID)
	categoryApp.Delete("/:categoryId", categoryHandler.DeleteCategory)

	// Content route
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Post("/create", contentHandler.CreateContent)
	contentApp.Get("/:contentId", contentHandler.GetContentByID)
	contentApp.Put("/:contentId", contentHandler.EditContentByID)
	contentApp.Delete("/:contentId", contentHandler.DeleteContent)
	contentApp.Post("/upload-image", contentHandler.UploadImageR2)

	// User route
	userApp := adminApp.Group("/users")
	userApp.Get("/profile", userHandler.GetUserById)
	userApp.Put("/update-password", userHandler.UpdatePassword)

	// FE
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentId", contentHandler.GetContentDetail)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	app.ShutdownWithContext(ctx)
}
