package ht

import (
	"context"
	"net/http"
	"simple_api/internal/service"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware для проверки токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем токен из заголовка Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Убираем префикс "Bearer " из токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Мы используем тот же секретный ключ для проверки
			return service.SecretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Передаем запрос дальше с данными из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Извлекаем user_id из токена
		userID := claims["user_id"].(float64)

		// Добавляем данные о пользователе в контекст
		ctx := context.WithValue(r.Context(), "user_id", int(userID))
		r = r.WithContext(ctx)

		// Переходим к следующему обработчику
		next.ServeHTTP(w, r)
	})
}