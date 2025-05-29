package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gokeeper/internal/config"
	"gokeeper/internal/handlers"
	"gokeeper/internal/middleware"
	"gokeeper/internal/repository/postgres"
	"gokeeper/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	cfg := config.New()

	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations("migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userService := service.NewUserService(db)
	secretService := service.NewSecretService(db, cfg.AESKey)

	userHandler := handlers.NewUserHandler(userService, cfg)
	secretHandler := handlers.NewSecretHandler(secretService)

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		secrets := api.Group("/secrets", middleware.JWTMiddleware(cfg))
		{
			secrets.POST("", secretHandler.Create)
			secrets.GET("", secretHandler.List)
			secrets.GET("/:id", secretHandler.Get)
			secrets.DELETE("/:id", secretHandler.Delete)
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"db_dsn":  cfg.GetDSN(),
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
