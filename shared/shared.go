package shared

import (
	"fmt"
	"strings"
)

const SAMPLE_SIZE = 10000
const SERVER_PORT = 46000
const AUTO = true

type Request struct {
	Player1 string
	Player2 string
}

type Reply struct {
	Result int
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ProcessaSolicitacao(request Request) int {
	possibilities := []string{"P", "A", "T"}

	if !inArray(request.Player1, possibilities) {
		return -1
	}
	if !inArray(request.Player2, possibilities) {
		return -1
	}
	if request.Player1 == request.Player2 {
		return 0
	}
	switch request.Player1 {
	case "P":
		if request.Player2 == "A" {
			return 2
		} else {
			return 1
		}
	case "A":
		if request.Player2 == "P" {
			return 1
		} else {
			return 2
		}
	case "T":
		if request.Player2 == "P" {
			return 2
		} else {
			return 1
		}
	default:
		return -1
	}
}

func GetMoves(auto bool) (player1Move string, player2Move string) {
	if auto {
		player1Move = "A"
		player2Move = "P"
	} else {
		fmt.Println("Favor informar a jogada do Player 1: (P = Pedra, A = Papel, T = Tesoura):")
		fmt.Print("\033[8m") // Hide input
		_, err := fmt.Scanln(&player1Move)
		if err != nil {
			panic(err)
		}
		fmt.Print("\033[28m") // Show input
		player1Move = strings.ToUpper(player1Move)

		fmt.Println("Favor informar a jogada do Player 2: (P = Pedra, A = Papel, T = Tesoura):")
		fmt.Print("\033[8m") // Hide input
		_, err = fmt.Scanln(&player2Move)
		if err != nil {
			panic(err)
		}
		fmt.Print("\033[28m") // Show input
		player2Move = strings.ToUpper(player2Move)
	}

	fmt.Println("Jogadas => Player 1:", player1Move, "Player 2:", player2Move)

	return player1Move, player2Move
}
