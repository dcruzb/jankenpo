package main

import (
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

	addr, err := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(shared.SERVER_PORT))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// connect to server
	conn, err := net.DialUDP("udp", nil, addr)
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

	var msgFromServer shared.Reply
	message := make([]byte, 4)

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		fmt.Println("Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		_, err := conn.Write([]byte(msgToServer + "\n"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// receive reply from server
		n, _, err := conn.ReadFromUDP(message)
		messageS := strings.TrimSuffix(string(message[:n]), "\n")
		result, err := strconv.Atoi(messageS)
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
