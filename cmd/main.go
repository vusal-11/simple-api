package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"simple_api/internal/transport/ht"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Используем драйвер для PostgreSQL
)

func main() {

	connStr := "postgres://postgres:post7027472@localhost/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Создаем маршруты
	r := mux.NewRouter()

	// Инициализируем обработчики
	ht.NewHandler(r, db)

	// Запускаем сервер
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
