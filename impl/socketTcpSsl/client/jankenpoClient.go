package client

import (
	"github.com/dcbCIn/jankenpo/impl/socketTcpSsl"
	"github.com/dcbCIn/jankenpo/shared"
	"os"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/socketTcpSsl/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string
	var sockTcpSsl socketTcpSsl.SocketTcpSsl

	// connect to server
	sockTcpSsl.ConnectToServer("localhost", strconv.Itoa(shared.TCP_SSL_PORT))

	shared.PrintlnInfo(NAME, "Connected successfully")
	shared.PrintlnInfo(NAME)

	// fecha o socket no final
	defer sockTcpSsl.CloseConnection()

	var msgFromServer shared.Reply

	// loop
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		start := time.Now()
		sockTcpSsl.Write(msgToServer)

		// receive reply from server
		message := sockTcpSsl.Read()
		elapsed += time.Since(start)
		message = strings.TrimSuffix(message, "\n")
		result, err := strconv.Atoi(message)
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
		time.Sleep(shared.WAIT * time.Millisecond)
	}
	return elapsed
}
