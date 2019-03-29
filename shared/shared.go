package shared

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/shared"

// Config
const CONECTIONS = 1
const TCP_PORT = 46000
const UDP_PORT = 47000
const JSON_PORT = 48000

// Debug
var SHOW_MESSAGES = []int{0, 1, 2}

const AUTO = true
const SAMPLE_SIZE = 3

type DebugLevel int

const (
	ERROR   DebugLevel = 0
	INFO    DebugLevel = 1
	MESSAGE DebugLevel = 2
)

func (d DebugLevel) ToInt() int {
	return [...]int{0, 1, 2}[d]
}

type Request struct {
	Player1 string
	Player2 string
}

type Reply struct {
	Result int
}

type Client struct {
	conn net.Conn
}

func Println(program string, messageLevel DebugLevel, message ...interface{}) {
	if len(SHOW_MESSAGES) > 0 {
		if inArrayInt(messageLevel.ToInt(), SHOW_MESSAGES) {
			switch messageLevel {
			case INFO:
				log.Println(program, "- INFO -", message)
			case MESSAGE:
				fmt.Println(message...)
			case ERROR:
				_, file, line, ok := runtime.Caller(1)
				if !ok {
					file = "???"
					line = 0
				}

				log.Println(program, "\n          ***** ERROR *****",
					"\n          File:", file,
					"\n          Line:", strconv.Itoa(line),
					"\n          Message:\n               ", message)
			}
		}
	}
}

func PrintlnInfo(program string, message ...interface{}) {
	Println(program, INFO, message...)
}

func PrintlnMessage(program string, message ...interface{}) {
	Println(program, MESSAGE, message...)
}

func PrintlnError(program string, message ...interface{}) {
	Println(program, ERROR, message...)
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func inArrayInt(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func randomMove() (move string) {
	//return "A" //For better performance dont return a random move

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	mv := rd.Intn(3)
	switch mv {
	case 0:
		move = "P"
	case 1:
		move = "A"
	case 2:
		move = "T"
	}

	return move
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
		player1Move = randomMove()
		player2Move = randomMove()
	} else {
		PrintlnMessage(NAME, "Favor informar a jogada do Player 1: (P = Pedra, A = Papel, T = Tesoura):")
		fmt.Print("\033[8m") // Hide input
		_, err := fmt.Scanln(&player1Move)
		if err != nil {
			panic(err)
		}
		fmt.Print("\033[28m") // Show input
		player1Move = strings.ToUpper(player1Move)

		PrintlnMessage(NAME, "Favor informar a jogada do Player 2: (P = Pedra, A = Papel, T = Tesoura):")
		fmt.Print("\033[8m") // Hide input
		_, err = fmt.Scanln(&player2Move)
		if err != nil {
			panic(err)
		}
		fmt.Print("\033[28m") // Show input
		player2Move = strings.ToUpper(player2Move)
	}

	PrintlnMessage(NAME, "Jogadas => Player 1:", player1Move, "Player 2:", player2Move)

	return player1Move, player2Move
}
