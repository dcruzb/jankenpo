package main

import (
	"encoding/json"
	"fmt"
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
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
	jsonDecoder := json.NewDecoder(conn)
	jsonEncoder := json.NewEncoder(conn)

	var msgFromServer shared.Reply

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		fmt.Println("Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := shared.Request{player1Move, player2Move}

		// send request to server
		err = jsonEncoder.Encode(msgToServer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// receive reply from server
		err = jsonDecoder.Decode(&msgFromServer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch msgFromServer.Result {
		case -1:
			fmt.Println("Invalid move")
		case 0:
			fmt.Println("Draw")
		default:
			fmt.Println("The winner is Player", msgFromServer.Result)
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
