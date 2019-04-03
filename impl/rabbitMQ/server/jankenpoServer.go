package server

import (
	"fmt"
	"jankenpo/impl/rabbitMQ"
	"jankenpo/shared"
	"os"
	"strconv"
	"sync"
)

const NAME = "jankenpo/rabbitMQ/server"

func waitForConection(rMQ rabbitMQ.RabbitMQ, idx int) {
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "started")

	// aceita conexões na porta
	//client := rMQ.WaitForConnection(idx)

	// fecha a conexão
	//defer client.CloseConnection()

	rMQ.CreateQueue("moves")
	rMQ.CreateQueue("result")

	var msgFromClient shared.Request

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (rabbitMQ)")
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

		// recebe solicitações do cliente
		message := rMQ.ReadOne("moves")

		shared.PrintlnInfo(NAME, "Message received: ", message)
		_, err := fmt.Sscanf(message, "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		rMQ.Write("result", strconv.Itoa(r))
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (rabbitMQ)")
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "ended")
}

func StartJankenpoServer() {
	var wg sync.WaitGroup
	shared.PrintlnInfo(NAME, "Initializing server rabbitMQ")

	// escuta na porta rabbitMQ configurada
	var rMQ rabbitMQ.RabbitMQ
	//rMQ.StartServer("", strconv.Itoa(shared.TCP_PORT), false, shared.CONECTIONS)
	rMQ.ConnectToServer("localhost", strconv.Itoa(shared.RABBITMQ_PORT))
	defer rMQ.CloseConnection()

	for idx := 0; idx < shared.CONECTIONS; idx++ {
		wg.Add(1)
		go func(i int) {
			waitForConection(rMQ, i)

			wg.Done()
		}(idx)
	}
	wg.Wait()
	shared.PrintlnInfo(NAME, "Fim do Servidor rabbitMQ")
}
