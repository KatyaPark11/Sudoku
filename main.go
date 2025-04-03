package main

import (
	"encoding/json"
	"net/http"
)

const N = 9

var board [N][N]int
var steps [][N][N]int // Список шагов решения

func isSafe(board [N][N]int, row, col, num int) bool {
	for x := range N {
		if board[row][x] == num || board[x][col] == num || board[row/3*3+x/3][col/3*3+x%3] == num {
			return false
		}
	}
	return true
}

func hiddenSingles(board *[N][N]int) bool {
	found := false

	// Проверяем каждое число от 1 до 9
	for num := 1; num <= 9; num++ {
		// Проверяем каждый бокс
		for boxRow := range 3 {
			for boxCol := range 3 {
				// Определяем координаты верхнего левого угла бокса
				startRow := boxRow * 3
				startCol := boxCol * 3

				// Список для хранения возможных позиций для текущего числа в боксе
				possiblePositions := []struct{ row, col int }{}

				// Ищем возможные позиции для числа в текущем боксе
				for i := range 3 {
					for j := range 3 {
						row := startRow + i
						col := startCol + j

						if board[row][col] == 0 && isSafe(*board, row, col, num) {
							possiblePositions = append(possiblePositions, struct{ row, col int }{row, col})
						}
					}
				}

				// Если только одна возможная позиция для числа в боксе
				if len(possiblePositions) == 1 {
					board[possiblePositions[0].row][possiblePositions[0].col] = num
					found = true
				}
			}
		}

		// Проверяем каждую строку и столбец
		for i := range N {
			possibleRowPos := -1
			possibleColPos := -1

			// Проверяем строки и столбцы на наличие единственного кандидата для числа
			for j := range N {
				if board[i][j] == 0 && isSafe(*board, i, j, num) {
					if possibleRowPos == -1 {
						possibleRowPos = j // Запоминаем первую возможную позицию
					} else {
						possibleRowPos = -2 // Больше одной позиции найдено
						break
					}
				}
			}

			if possibleRowPos != -2 && possibleRowPos != -1 { // Если нашли единственную позицию в строке
				board[i][possibleRowPos] = num
				found = true
			}

			for j := range N {
				if board[j][i] == 0 && isSafe(*board, j, i, num) {
					if possibleColPos == -1 {
						possibleColPos = j // Запоминаем первую возможную позицию
					} else {
						possibleColPos = -2 // Больше одной позиции найдено
						break
					}
				}
			}

			if possibleColPos != -2 && possibleColPos != -1 { // Если нашли единственную позицию в столбце
				board[possibleColPos][i] = num
				found = true
			}
		}
	}

	return found // Возвращаем true если были сделаны изменения на доске.
}

func solveSudoku(board *[N][N]int, isSteps bool) bool {
	strategies := []func(*[N][N]int) bool{
		hiddenSingles,

		// Добавьте остальные стратегии в порядке их приоритета...
	}

	for {
		madeProgress := false

		for _, strategy := range strategies {
			if strategy(board) {
				madeProgress = true
				if isSteps {
					saveStep(*board) // Сохраняем текущий шаг после применения стратегии
				}
			}
		}

		if !madeProgress { // Если не было сделано прогресса, значит, нужно использовать пробу и ошибку
			break
		}
	}

	return backtrackSolve(board, isSteps)
}

func backtrackSolve(board *[N][N]int, isSteps bool) bool {
	for row := range N {
		for col := range N {
			if board[row][col] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafe(*board, row, col, num) {
						board[row][col] = num
						if isSteps {
							saveStep(*board) // Сохраняем текущий шаг
						}
						if backtrackSolve(board, isSteps) {
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
