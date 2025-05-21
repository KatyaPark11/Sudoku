package main

import (
	context "context"

	pb "C:/Users/KatyaPark11/Documents/GitHub/Sudoku/gen/sudoku.pb.go"
)

// SudokuService реализует интерфейс pb.SudokuServiceServer
type SudokuService struct {
	pb.UnimplementedSudokuServiceServer
}

// Решает судоку (заглушка)
func (s *SudokuService) Solve(ctx context.Context, req *pb.SolveRequest) (*pb.SolveResponse, error) {
	puzzle := req.GetPuzzle()

	// Тут должна быть логика решения судоку.
	// Для примера возвращаем тот же массив.

	solution := make([]int32, len(puzzle))
	copy(solution, puzzle)

	// В реальности нужно реализовать алгоритм решения.

	return &pb.SolveResponse{
		Solution: solution,
	}, nil
}
