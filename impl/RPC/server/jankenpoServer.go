package server

import (
	"github.com/dcbCIn/jankenpo/impl/RPC"
	"github.com/dcbCIn/jankenpo/shared"
	"strconv"
	"sync"
)

const NAME = "jankenpo/rpc/server"

func waitForConection(rpc rpc.RPC, idx int) {
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "started")

	// fecha o socket
	defer rpc.CloseConnection()

	// aceita conexões na porta
	rpc.WaitForConnection(idx)

	shared.PrintlnInfo(NAME, "Servidor finalizado (RPC)")
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "ended")
}

func StartJankenpoServer() {
	var wg sync.WaitGroup
	shared.PrintlnInfo(NAME, "Initializing server RPC")

	// escuta na porta tcp configurada
	var rpc rpc.RPC
	rpc.StartServer("", strconv.Itoa(shared.RPC_PORT))
	defer rpc.StopServer()

	for idx := 0; idx < shared.CONECTIONS; idx++ {
		wg.Add(1)
		go func(i int) {
			waitForConection(rpc, i)

			wg.Done()
		}(idx)
	}
	wg.Wait()
	shared.PrintlnInfo(NAME, "Fim do Servidor RPC")
}
