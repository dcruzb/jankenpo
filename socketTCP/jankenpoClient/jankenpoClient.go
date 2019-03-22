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

	// fecha o socket no final
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	// create a decoder/encoder
	var msgFromServer shared.Reply

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		fmt.Println("Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		_, err := fmt.Fprintf(conn, msgToServer+"\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// receive reply from server
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
		PlayJanKenPo(shared.AUTO)
		wg.Done()
	}()

	wg.Wait()
}
