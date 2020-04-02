package client

import (
	"encoding/json"
	"github.com/dcbCIn/jankenpo/shared"
	"net"
	"os"
	"strconv"
	"time"
)

const NAME = "jankenpo/socketJson/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string

	// connect to server
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(shared.JSON_PORT))
	if err != nil {
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}
	shared.PrintlnInfo(NAME, "Connected successfully")
	shared.PrintlnInfo(NAME)

	// create a decoder/encoder
	jsonDecoder := json.NewDecoder(conn)
	jsonEncoder := json.NewEncoder(conn)

	var msgFromServer shared.Reply

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := shared.Request{player1Move, player2Move}

		// send request to server
		err = jsonEncoder.Encode(msgToServer)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// receive reply from server
		err = jsonDecoder.Decode(&msgFromServer)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		switch msgFromServer.Result {
		case -1:
			shared.PrintlnMessage(NAME, "Invalid move")
		case 0:
			shared.PrintlnMessage(NAME, "Draw")
		default:
			shared.PrintlnMessage(NAME, "The winner is Player", msgFromServer.Result)
		}
		shared.PrintlnMessage(NAME, "------------------------------------------------------------------")
		shared.PrintlnMessage(NAME)
	}
	elapsed = time.Since(start)
	return elapsed
}
