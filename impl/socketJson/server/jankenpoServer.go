package server

import (
	"encoding/json"
	"github.com/dcbCIn/jankenpo/shared"
	"net"
	"os"
	"strconv"
)

const NAME = "jankenpo/socketJson/server"

func StartJankenpoServer() {
	shared.PrintlnInfo(NAME, "Initializing server Json")
	// escuta na porta tcp configurada
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(shared.JSON_PORT))

	// aceita conexões na porta
	conn, err := ln.Accept()
	if err != nil {
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}

	// fecha o socket
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	// cria um cofificador/decodificador Json
	jsonDecoder := json.NewDecoder(conn)
	jsonEncoder := json.NewEncoder(conn)
	var msgFromClient shared.Request

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (Json over TCP)...")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente e decodifica-as
		err = jsonDecoder.Decode(&msgFromClient)
		if err != nil {
			shared.PrintlnError(NAME, err)
			return
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		msgToClient := shared.Reply{r}
		err = jsonEncoder.Encode(msgToClient)
		if err != nil {
			shared.PrintlnError(NAME, err)
			return
		}
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (Json over TCP)")
}
