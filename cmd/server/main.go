package main

import (
	"chirik/internal/handlers"
	"chirik/internal/middleware"
	"chirik/internal/realtime"
	"chirik/internal/repository"
	"chirik/internal/router"
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Getenv("CHIRIK_JWT_SECRET")) < 32 {
		log.Fatal("CHIRIK_JWT_SECRET must contain at least 32 characters")
	}

	dataDir := envOr("CHIRIK_DATA_DIR", ".")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Fatal("Data directory error:", err)
	}
	dbPath := envOr("CHIRIK_DB_PATH", filepath.Join(dataDir, "chirik.db"))
	webDir := envOr("CHIRIK_WEB_DIR", "./web")

	// База данных хранится отдельно от каталога приложения, чтобы обновления не затронули данные.
	repo, err := repository.New(dbPath)
	if err != nil {
		log.Fatal("Database error:", err)
	}
	defer repo.Close()
	realtimePublisher := realtime.NewHTTPPublisher(os.Getenv("CHIRIK_REALTIME_URL"), os.Getenv("CHIRIK_REALTIME_SECRET"))

	// Handlers
	authHandler := handlers.NewAuthHandler(repo)
	tweetHandler := handlers.NewTweetHandler(repo, realtimePublisher)
	followHandler := handlers.NewFollowHandler(repo, realtimePublisher)
	messageHandler := handlers.NewMessageHandler(repo, realtimePublisher)
	usersHandler := handlers.NewUsersHandler(repo)

	// Router (API)
	r := router.Setup(authHandler, tweetHandler, followHandler, messageHandler, usersHandler)

	// Middleware для API
	apiHandler := middleware.CORS(r)

	// Основной мультиплексор
	mux := http.NewServeMux()

	// Один внешний порт: realtime остаётся внутренним и доступен через /events.
	realtimeURL, err := url.Parse(envOr("CHIRIK_REALTIME_PROXY_URL", "http://127.0.0.1:8090"))
	if err != nil {
		log.Fatal("Realtime proxy URL error:", err)
	}
	mux.Handle("/events", httputil.NewSingleHostReverseProxy(realtimeURL))

	// API маршруты
	mux.Handle("/api/", apiHandler)

	// Статические файлы (React build)
	fileServer := http.FileServer(http.Dir(webDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрос к API — пропускаем
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Проверяем, существует ли файл
		filePath := filepath.Join(webDir, filepath.Clean("/"+r.URL.Path))
		if _, err := os.Stat(filePath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Для всех остальных путей — отдаём index.html (React Router)
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	// Сервер
	srv := &http.Server{
		Addr:    envOr("CHIRIK_SERVER_ADDR", ":8080"),
		Handler: mux,
	}

	log.Printf("🚀 Чирик запущен на %s", srv.Addr)
	log.Println("📱 Открой в браузере по IP телефона и порту сервера")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	<-stop
	log.Println("⏳ Останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown error:", err)
	}

	log.Println("✅ Сервер остановлен")
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
