package client

import (
	"log"

	pbAuth "github.com/KatyaPark11/Sudoku/generated/auth"     // замените на ваш путь сгенерированных proto
	pbSudoku "github.com/KatyaPark11/Sudoku/generated/sudoku" // замените на ваш путь сгенерированных proto

	"google.golang.org/grpc"
)

var (
	authConn     *grpc.ClientConn
	sudokuConn   *grpc.ClientConn
	authClient   pbAuth.AuthServiceClient
	sudokuClient pbSudoku.SudokuServiceClient
)

func init() {
	var err error

	authConn, err = grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	authClient = pbAuth.NewAuthServiceClient(authConn)

	sudokuConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to sudoku service: %v", err)
	}
	sudokuClient = pbSudoku.NewSudokuServiceClient(sudokuConn)
}

func AuthClient() pbAuth.AuthServiceClient {
	return authClient
}

func SudokuClient() pbSudoku.SudokuServiceClient {
	return sudokuClient
}
