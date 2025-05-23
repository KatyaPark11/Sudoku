package main

import (
	"log"
	"net/http"

	pbAuth "github.com/KatyaPark11/Sudoku/generated/auth"     // замените на ваш путь сгенерированных proto
	pbSudoku "github.com/KatyaPark11/Sudoku/generated/sudoku" // замените на ваш путь сгенерированных proto

	"google.golang.org/grpc"
)

var (
	authClient   pbAuth.AuthServiceClient
	sudokuClient pbSudoku.SudokuServiceClient
)

func main() {
	var err error

	// Инициализация соединений с gRPC-сервисами
	authConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	authClient = pbAuth.NewAuthServiceClient(authConn)

	sudokuConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to sudoku service: %v", err)
	}
	sudokuClient = pbSudoku.NewSudokuServiceClient(sudokuConn)

	// Запуск HTTP-сервера или другого основного кода
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/solve", handleSolve)

	// Статические файлы (HTML)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
