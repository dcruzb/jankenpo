package main

import (
	"bufio"
	"fmt"
	"jankenpo/shared"
	"net"
	"os"
	"strconv"
	"sync"
)

func StartJankenpoServer() {
	fmt.Println("Initializing server")
	// escuta na porta tcp configurada
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(shared.SERVER_PORT))

	// aceita conexões na porta
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
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

	fmt.Println("Servidor pronto para receber solicitações (TCP)...")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println("Message received: ", message)
		_, err = fmt.Sscanf(message, "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		_, err = conn.Write([]byte(strconv.Itoa(r) + "\n"))
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
