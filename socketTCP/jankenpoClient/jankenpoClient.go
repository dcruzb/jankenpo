package main

import (
	"bufio"
	"fmt"
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func PlayJanKenPo(auto bool) {
	var player1Move, player2Move string

	// connect to server
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(shared.SERVER_PORT))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Connected successfully")
	fmt.Println()

	// create a decoder/encoder
	//jsonDecoder := json.NewDecoder(conn)
	//jsonEncoder := json.NewEncoder(conn)
	var msgFromServer shared.Reply

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		fmt.Println("Game", i)
		if auto {
			player1Move = "A"
			player2Move = "P"
		} else {
			fmt.Println("Favor informar a jogada do Player 1: (P = Pedra, A = Papel, T = Tesoura):")
			fmt.Print("\033[8m") // Hide input
			_, err := fmt.Scanln(&player1Move)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Print("\033[28m") // Show input
			player1Move = strings.ToUpper(player1Move)

			fmt.Println("Favor informar a jogada do Player 2: (P = Pedra, A = Papel, T = Tesoura):")
			fmt.Print("\033[8m") // Hide input
			_, err = fmt.Scanln(&player2Move)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Print("\033[28m") // Show input
			player2Move = strings.ToUpper(player2Move)
		}

		fmt.Println("Jogadas => Player 1: ", string(player1Move), "Player 2:", string(player2Move))

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		//err = jsonEncoder.Encode(msgToServer)
		// send to socket
		_, err := fmt.Fprintf(conn, msgToServer+"\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// receive reply from server
		//err = jsonDecoder.Decode(&msgFromServer)
		message, err := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		result, err := strconv.Atoi(message)
		if err != nil {
			fmt.Println(err)
			result = -1
		}
		msgFromServer = shared.Reply{result}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println()
		switch msgFromServer.Result {
		case 1, 2:
			fmt.Println("The winner is Player", msgFromServer.Result)
		case 0:
			fmt.Println("Draw")
		default:
			fmt.Println("Invalid move")
		}
		fmt.Println("------------------------------------------------------------------")
		fmt.Println()
	}
	elapsed := time.Since(start)
	fmt.Printf("Tempo: %s \n", elapsed)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		PlayJanKenPo(false)
		wg.Done()
	}()

	wg.Wait()
}
