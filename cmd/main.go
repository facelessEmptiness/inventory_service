// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"github.com/facelessEmptiness/inventory_service/internal/config"
	"github.com/facelessEmptiness/inventory_service/internal/infrastructure/repository/db"
	db2 "github.com/facelessEmptiness/inventory_service/internal/repository/db"

	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Инициализируем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к MongoDB
	mongodb, err := db.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Контекст для закрытия соединения при завершении
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer mongodb.Close(ctx)

	// Запускаем миграции MongoDB
	if err := db2.RunMigrations(mongodb.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Создаем экземпляры репозиториев
	productRepo := repository.NewProductRepository(mongodb.DB)
	categoryRepo := repository.NewCategoryRepository(mongodb.DB)

	// Создаем экземпляры use cases
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	// Создаем экземпляры middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	grpc.NewServer()

	// Создаем экземпляры обработчиков
	productHandler := handler.NewProductHandler(productUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)

	// Настраиваем Gin
	router := gin.Default()

	// Настраиваем маршруты
	handler.SetupRoutes(router, productHandler, categoryHandler, authMiddleware)

	// Запускаем сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting server on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ждем сигнала для корректного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Даем серверу 5 секунд на завершение текущих запросов
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
