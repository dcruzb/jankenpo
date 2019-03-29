package server

import (
	"bufio"
	"fmt"
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
	"sync"
)

const NAME = "jankenpo/socketTCP/server"

func waitForConection(ln net.Listener, idx int) {
	//var id int  = 1 //rand.Seed(now).Int(1000)
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "started")

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

	var msgFromClient shared.Request

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (TCP)")
	for i := 0; i < shared.SAMPLE_SIZE; i++ {

		// recebe solicitações do cliente
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}
		//shared.PrintlnInfo(NAME, "Message received: ", message)
		_, err = fmt.Sscanf(message, "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		_, err = conn.Write([]byte(strconv.Itoa(r) + "\n"))
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (TCP)")
	shared.PrintlnInfo(NAME, "Connection", strconv.Itoa(idx), "ended")
}

func StartJankenpoServer() {
	var wg sync.WaitGroup
	shared.PrintlnInfo(NAME, "Initializing server TCP")
	// escuta na porta tcp configurada
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(shared.TCP_PORT))
	defer func() {
		err := ln.Close()
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}
	}()

	for idx := 0; idx < shared.CONECTIONS; idx++ {
		wg.Add(1)
		go func(i int) {
			waitForConection(ln, i)

			wg.Done()
		}(idx)
	}
	wg.Wait()
	shared.PrintlnInfo(NAME, "Fim do Servidor TCP")
}
