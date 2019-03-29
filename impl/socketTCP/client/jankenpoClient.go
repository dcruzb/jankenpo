package client

import (
	"bufio"
	"io"
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const NAME = "jankenpo/socketTCP/client"

func PlayJanKenPo(auto bool) (elapsed time.Duration) {
	var player1Move, player2Move string

	// connect to server
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(shared.TCP_PORT))
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

	// create a decoder/encoder
	var msgFromServer shared.Reply

	// loop
	start := time.Now()
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		shared.PrintlnMessage(NAME, "Game", i)

		player1Move, player2Move = shared.GetMoves(auto)

		// prepare request
		msgToServer := player1Move + " " + player2Move //shared.Request{player1Move, player2Move}

		// send request to server
		// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
		//_, err := fmt.Fprintf(conn, msgToServer+"\n")
		//_, err := conn.Write([]byte( msgToServer + "\n"))
		/*reader := bufio.NewWriter(conn)
		_, err := reader.WriteString( msgToServer + "\n")
		reader.Flush()*/
		/*reader := bufio.NewWriter(conn)
		_, err := io.WriteString(reader, msgToServer + "\n")
		reader.Flush()*/
		_, err := io.WriteString(conn, msgToServer+"\n")
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// receive reply from server
		message, err := bufio.NewReader(conn).ReadString('\n')
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
	}
	elapsed = time.Since(start)
	return elapsed
}
