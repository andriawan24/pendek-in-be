package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/middlewares"
	"github.com/andriawan24/link-short/internal/routes"
	"github.com/andriawan24/link-short/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mostafa-asg/ip2country"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loadEnv()
	loadIPDatabase()

	rdb := newRedisClient()
	db := newPostgresDB(ctx)
	defer db.Close()

	queries := database.New(db)
	router := setupRouter(ctx, db, queries, rdb)
	server := newHTTPServer(router)

	gracefulShutdown(ctx, server)
	startServer(server)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to read .env file: %v", err)
	}
}

func loadIPDatabase() {
	if err := ip2country.Load("internal/sources/dbip-country.csv"); err != nil {
		log.Fatalf("Failed to load IP country database file: %v", err)
	}
}

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     getenv("REDIS_ADDR", "localhost:6379"),
		Password: getenv("REDIS_PASSWORD", ""),
		DB:       0,
		Protocol: 2,
	})
}

func newPostgresDB(ctx context.Context) *sql.DB {
	db, err := sql.Open("postgres", buildConnectionString())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

func buildConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", "postgres"),
		getenv("DB_NAME", "link-short"),
		getenv("DB_SSLMODE", "disable"),
	)
}

func setupRouter(ctx context.Context, db *sql.DB, queries *database.Queries, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)

	r.Use(cors.New(buildCORSConfig()))

	registerRoutes(r, ctx, db, queries, rdb)

	return r
}

func buildCORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     parseAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: getenv("CORS_ALLOW_CREDENTIALS", "true") == "true",
		MaxAge:           12 * time.Hour,
	}
}

func parseAllowedOrigins() []string {
	corsOrigins := getenv("CORS_ALLOW_ORIGINS", "*")
	if corsOrigins == "*" {
		return []string{"*"}
	}

	var origins []string
	for origin := range strings.SplitSeq(corsOrigins, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func registerRoutes(r *gin.Engine, ctx context.Context, db *sql.DB, queries *database.Queries, rdb *redis.Client) {
	userService := services.NewUserService(ctx, queries)
	linkService := services.NewLinkService(ctx, queries)
	cacheService := services.NewCacheService(rdb)
	clickLogService := services.NewClickLogService(queries)

	linkRoutes := routes.NewLinkRoutes(linkService, clickLogService, cacheService)
	authRoutes := routes.NewAuthRoutes(userService)
	analyticRoutes := routes.NewAnalyticRoutes(linkService, clickLogService)

	authGroup := r.Group("/auth")
	{
		authGroup.GET("/me", middlewares.RequiredAuth(), authRoutes.Profile)
		authGroup.POST("/login", authRoutes.Login)
		authGroup.POST("/refresh", authRoutes.Refresh)
		authGroup.POST("/register", authRoutes.Register)
		authGroup.PUT("/update-profile", middlewares.RequiredAuth(), authRoutes.UpdateProfile)
	}

	linkGroup := r.Group("/links", middlewares.RequiredAuth())
	{
		linkGroup.GET("/all", linkRoutes.GetLinks)
		linkGroup.GET("/:id", linkRoutes.GetLink)
		linkGroup.POST("/create", linkRoutes.InsertLink)
	}

	analyticGroup := r.Group("/analytics", middlewares.RequiredAuth())
	{
		analyticGroup.GET("/dashboard", analyticRoutes.GetDashboard)
		analyticGroup.GET("/", analyticRoutes.GetAnalytics)
	}

	r.GET("/:code", linkRoutes.Redirect)
	r.GET("/health", healthCheckHandler(db))
	r.NoRoute()
}

func healthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthCtx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		if err := db.PingContext(healthCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "db unavailable"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func newHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + getenv("HTTP_PORT", "8080"),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func gracefulShutdown(ctx context.Context, srv *http.Server) {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()
}

func startServer(srv *http.Server) {
	log.Printf("HTTP server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
