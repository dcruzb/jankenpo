package server

import (
	"fmt"
	"github.com/dcbCIn/jankenpo/impl/quic"
	"github.com/dcbCIn/jankenpo/shared"
	"os"
	"strconv"
	"sync"
)

const NAME = "jankenpo/quic/server"

func waitForConection(quic quic.Quic, idx int) {
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "started")

	// aceita conexões na porta
	client := quic.WaitForConnection(idx)

	// fecha o socket
	//defer client.CloseConnection()

	var msgFromClient shared.Request

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (Quic)")
	for i := 0; i <= shared.SAMPLE_SIZE; i++ {
		// recebe solicitações do cliente
		message := client.Read()

		shared.PrintlnInfo(NAME, "Message received: ", message)
		_, err := fmt.Sscanf(message, "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		client.Write(strconv.Itoa(r) + "\n")
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (Quic)")
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "ended")
}

func StartJankenpoServer() {
	var wg sync.WaitGroup
	shared.PrintlnInfo(NAME, "Initializing server Quic")

	// escuta na porta quic configurada
	var quic quic.Quic
	quic.StartServer("", strconv.Itoa(shared.QUIC_PORT), false, shared.CONECTIONS)
	defer quic.StopServer()

	for idx := 0; idx < shared.CONECTIONS; idx++ {
		wg.Add(1)
		go func(i int) {
			waitForConection(quic, i)

			wg.Done()
		}(idx)
	}
	wg.Wait()
	shared.PrintlnInfo(NAME, "Fim do Servidor Quic")
}
