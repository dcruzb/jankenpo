package shared

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/shared"

// Config
const NAME_SERVER_IP = "127.0.0.1"
const NAME_SERVER_PORT = 45000
const QUIC_PORT = 51000
const TCP_PORT = 46000
const UDP_PORT = 47000
const JSON_PORT = 48000
const RPC_PORT = 49000
const RABBITMQ_PORT = 5672
const MID_PORT = 50000
const CONECTIONS = 1

// Debug
const AUTO = true
const SAMPLE_SIZE = 10000
const QUIC = true
const SOCKET_TCP = true
const SOCKET_UDP = false
const JSON = false
const RPC = false
const RABBIT_MQ = false
const MID = false
const WAIT = 5 // tempo em ms

var SHOW_MESSAGES = []DebugLevel{ERROR} //, INFO, MESSAGE}

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

func (rq Request) Play(request Request, reply *Reply) error {
	*reply = Reply{ProcessaSolicitacao(request)}
	return nil
}

func Println(program string, messageLevel DebugLevel, message ...interface{}) {
	if len(SHOW_MESSAGES) > 0 {
		if inArrayDL(messageLevel, SHOW_MESSAGES) {
			switch messageLevel {
			case INFO:
				var logs []interface{}
				logs = append(logs, program, "- INFO -")
				logs = append(logs, message...)
				log.Println(logs...)
			case MESSAGE:
				fmt.Println(message...)
			case ERROR:
				_, file, line, ok := runtime.Caller(2)
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

func FailOnError(program string, err error, msg string) {
	if err != nil {
		Println(program, ERROR, msg, ":", err)
		os.Exit(1)
	}
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func inArrayDL(a DebugLevel, list []DebugLevel) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func randomMove() (move string) {
	return "A" //For better performance dont return a random move

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

// Setup a bare-bones TLS config for the server
func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(crand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(crand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"exemplo"},
	}
}