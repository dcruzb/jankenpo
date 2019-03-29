package socketUDP

import (
	"fmt"
	"jankenpo/shared"
	"log"
	"net"
	"os"
	"strconv"
)

func StartJankenpoServer() {
	shared.PrintlnInfo(NAME, "Initializing server UDP")

	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(shared.UDP_PORT))
	if err != nil {
		log.Fatal(err)
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}

	// escuta na porta tcp configurada
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		shared.PrintlnError(NAME, err)
		os.Exit(1)
	}

	// fecha o socket no final
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	// cria um cofificador/decodificador Json
	var msgFromClient shared.Request
	message := make([]byte, 4)

	shared.PrintlnInfo(NAME, "Servidor pronto para receber solicitações (UDP)")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente
		n, addr, err := conn.ReadFromUDP(message)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}
		//shared.PrintlnInfo(NAME, "Message received: ", message)
		_, err = fmt.Sscanf(string(message[:n]), "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		_, err = conn.WriteTo([]byte(strconv.Itoa(r)+"\n"), addr)
		if err != nil {
			shared.PrintlnError(NAME, err)
			os.Exit(1)
		}
	}
	shared.PrintlnInfo(NAME, "Servidor finalizado (UDP)")
}
