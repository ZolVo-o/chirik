package middleware

import (
    "context"
    "net/http"
    "strings"
    "chirik/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Получаем токен из заголовка
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, `{"error":"Authorization required"}`, http.StatusUnauthorized)
            return
        }

        // Проверяем формат
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, `{"error":"Invalid authorization format"}`, http.StatusUnauthorized)
            return
        }

        // Валидируем токен
        claims, err := auth.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, `{"error":"Invalid token: `+err.Error()+`"}`, http.StatusUnauthorized)
            return
        }

        // Сохраняем user_id в контексте
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next(w, r.WithContext(ctx))
    }
}

// Helper для получения user_id из контекста
func GetUserID(r *http.Request) (int, bool) {
    userID, ok := r.Context().Value(UserIDKey).(int)
    return userID, ok
}
