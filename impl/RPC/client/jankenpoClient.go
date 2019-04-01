package client

import (
	"jankenpo/impl/RPC"
	"jankenpo/shared"
	"strconv"
	"time"
)

const NAME = "jankenpo/rpc/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string
	var rpc rpc.RPC

	// connect to server
	rpc.ConnectToServer("localhost", strconv.Itoa(shared.RPC_PORT))

	shared.PrintlnInfo(NAME, "Connected successfully")
	shared.PrintlnInfo(NAME)

	// Close the connection
	defer rpc.CloseConnection()

	var msgFromServer *shared.Reply
	var msgToServer shared.Request

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer = shared.Request{player1Move, player2Move} //player1Move + " " + player2Move //

		// send request to server and receive reply at the same time
		msgFromServer = rpc.Call("Request.Play", msgToServer)

		shared.PrintlnMessage(NAME)
		switch (*msgFromServer).Result {
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
