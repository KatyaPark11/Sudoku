package main

import (
	"encoding/json"
	"net/http"
)

const N = 9

var board [N][N]int
var steps [][N][N]int // Список шагов решения

func isSafe(board [N][N]int, row, col, num int) bool {
	for x := 0; x < N; x++ {
		if board[row][x] == num || board[x][col] == num || board[row/3*3+x/3][col/3*3+x%3] == num {
			return false
		}
	}
	return true
}

func solveSudoku(board *[N][N]int, isSteps bool) bool {
	for row := 0; row < N; row++ {
		for col := 0; col < N; col++ {
			if board[row][col] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafe(*board, row, col, num) {
						board[row][col] = num
						if isSteps {
							saveStep(*board) // Сохраняем текущий шаг
						}
						if solveSudoku(board, isSteps) {
							return true
						}
						board[row][col] = 0
					}
				}
				return false
			}
		}
	}
	return true
}

func saveStep(currentBoard [N][N]int) {
	stepCopy := currentBoard        // Копируем текущее состояние доски
	steps = append(steps, stepCopy) // Добавляем в список шагов
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input [N][N]int
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		steps = nil
		board = input

		if solveSudoku(&board, false) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(board)
		} else {
			http.Error(w, "No solution exists", http.StatusBadRequest)
		}

		return
	}
}

func stepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input [N][N]int
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		steps = nil
		board = input

		if solveSudoku(&board, true) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(steps)
		} else {
			http.Error(w, "No solution exists", http.StatusBadRequest)
		}

		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/solve", solveHandler)
	http.HandleFunc("/step", stepHandler)
	http.ListenAndServe(":8080", nil)
}
