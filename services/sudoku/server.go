package main

import (
	"context"
	"log"
	"net"

	pb "github.com/KatyaPark11/Sudoku/generated/sudoku"
	"google.golang.org/grpc"
)

type sudokuServer struct {
	pb.UnimplementedSudokuServiceServer
}

// Простая функция-решатель судоку (заглушка)
func solveSudoku(puzzle string) string {
	// Тут должна быть логика решения судоку.
	// Для примера возвращаем тот же пазл.
	return puzzle // или можно вернуть фиксированный ответ для теста.
}

func (s *sudokuServer) Solve(ctx context.Context, req *pb.SudokuRequest) (*pb.SudokuResponse, error) {
	solution := solveSudoku(req.Puzzle)
	return &pb.SudokuResponse{Solution: solution}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSudokuServiceServer(grpcServer, &sudokuServer{})

	log.Println("Sudoku service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
