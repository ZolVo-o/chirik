package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
    "chirik/internal/handlers"
    "chirik/internal/middleware"
    "chirik/internal/repository"
    "chirik/internal/router"
)

func main() {
    // База данных
    repo, err := repository.New("./chirik.db")
    if err != nil {
        log.Fatal("Database error:", err)
    }
    defer repo.Close()

    // Handlers
    authHandler := handlers.NewAuthHandler(repo)
    tweetHandler := handlers.NewTweetHandler(repo)
    followHandler := handlers.NewFollowHandler(repo)

    // Router (API)
    r := router.Setup(authHandler, tweetHandler, followHandler)

    // Middleware для API
    apiHandler := middleware.CORS(r)

    // Основной мультиплексор
    mux := http.NewServeMux()

    // API маршруты
    mux.Handle("/api/", apiHandler)

    // Статические файлы (React build)
    fileServer := http.FileServer(http.Dir("./web"))
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Если запрос к API — пропускаем
        if strings.HasPrefix(r.URL.Path, "/api/") {
            http.NotFound(w, r)
            return
        }

        // Проверяем, существует ли файл
        filePath := "./web" + r.URL.Path
        if _, err := os.Stat(filePath); err == nil {
            fileServer.ServeHTTP(w, r)
            return
        }

        // Для всех остальных путей — отдаём index.html (React Router)
        http.ServeFile(w, r, "./web/index.html")
    })

    // Сервер
    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    log.Println("🚀 Чирик запущен на http://localhost:8080")
    log.Println("📱 Открой в браузере: http://localhost:8080")

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
