package main

import (
	"fmt"
	"jankenpo/shared"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

func StartJankenpoServer() {
	fmt.Println("Initializing server")

	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(shared.SERVER_PORT))
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		os.Exit(1)
	}

	// escuta na porta tcp configurada
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
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

	fmt.Println("Servidor pronto para receber solicitações (UDP)...")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente
		n, addr, err := conn.ReadFromUDP(message)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println("Message received: ", message)
		_, err = fmt.Sscanf(string(message[:n]), "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		_, err = conn.WriteTo([]byte(strconv.Itoa(r)+"\n"), addr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		StartJankenpoServer()
		wg.Done()
	}()

	wg.Wait()
}
