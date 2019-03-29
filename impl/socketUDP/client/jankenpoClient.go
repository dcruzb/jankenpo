package client

import (
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/socketUDP/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string

	addr, err := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(shared.UDP_PORT))
	if err != nil {
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}

	// connect to server
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}
	shared.PrintlnInfo(NAME, "Connected successfully")
	shared.PrintlnInfo(NAME)

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
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		_, err := conn.Write([]byte(msgToServer + "\n"))
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// receive reply from server
		n, _, err := conn.ReadFromUDP(message)
		messageS := strings.TrimSuffix(string(message[:n]), "\n")
		result, err := strconv.Atoi(messageS)
		if err != nil {
			shared.PrintlnError(NAME, err)
			result = -1
		}
		msgFromServer = shared.Reply{result}
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		shared.PrintlnMessage(NAME)
		switch msgFromServer.Result {
		case 1, 2:
			shared.PrintlnMessage(NAME, "The winner is Player", msgFromServer.Result)
		case 0:
			shared.PrintlnMessage(NAME, "Draw")
		default:
			shared.PrintlnMessage(NAME, "Invalid move")
		}
		shared.PrintlnMessage(NAME, "------------------------------------------------------------------")
		shared.PrintlnMessage(NAME)
	}
	elapsed = time.Since(start)
	return elapsed
}
