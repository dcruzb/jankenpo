package client

import (
	"fmt"
	"github.com/dcbCIn/jankenpo/impl/quic"
	"github.com/dcbCIn/jankenpo/shared"
	"os"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/quic/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string
	var quic quic.Quic
	filename := fmt.Sprintf("%s%s%s", "./quic", time.Now(), ".csv")

	// connect to server
	quic.ConnectToServer("localhost", strconv.Itoa(shared.QUIC_PORT))

	shared.PrintlnInfo(NAME, "Connected successfully")
	shared.PrintlnInfo(NAME)

	// fecha o socket no final
	//defer quic.CloseConnection()

	var msgFromServer shared.Reply

	shared.WriteToFile(filename, "Type;SAMPLE_SIZE;WAIT;InvNumber;unitaryElapsed;unitaryElapsedNanoseconds\n")
	// loop
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		start := time.Now()
		quic.Write(msgToServer)

		// receive reply from server
		message := quic.Read()
		unitaryElapsed := time.Since(start)
		shared.WriteToFile(filename, fmt.Sprintf("%s%d%s%d%s%d%s%s%s%d\n","Quic;",shared.SAMPLE_SIZE, ";", shared.WAIT, ";", i, ";", unitaryElapsed, ";", unitaryElapsed.Nanoseconds()))
		elapsed += unitaryElapsed

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

	shared.WriteToFile(filename, fmt.Sprintf("%s%d%s%d%s%d%s%s\n","TCP/Ssl;",shared.SAMPLE_SIZE, ";", shared.WAIT, ";", elapsed, ";", elapsed / shared.SAMPLE_SIZE))
	return elapsed
}
