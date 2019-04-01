package server

import (
	"fmt"
	"jankenpo/impl/RPC"
	"jankenpo/shared"
	"os"
	"strconv"
	"sync"
)

const NAME = "jankenpo/rpc/server"

func waitForConection(rpc rpc.RPC, idx int) {
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "started")

	// aceita conexões na porta
	client := rpc.WaitForConnection(idx)

	// fecha o socket
	defer client.CloseConnection()

	var msgFromClient shared.Request

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (RPC)")
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

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
		client.Write(strconv.Itoa(r))
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (RPC)")
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "ended")
}

func StartJankenpoServer() {
	var wg sync.WaitGroup
	shared.PrintlnInfo(NAME, "Initializing server RPC")

	// escuta na porta tcp configurada
	var rpc rpc.RPC
	rpc.StartServer("", strconv.Itoa(shared.RPC_PORT), false, shared.CONECTIONS)
	defer rpc.StopServer()

	/*for idx := 0; idx < shared.CONECTIONS; idx++ {
		wg.Add(1)
		go func(i int) {
			waitForConection(rpc, i)

			wg.Done()
		}(idx)
	}*/
	wg.Wait()
	shared.PrintlnInfo(NAME, "Fim do Servidor RPC")
}
