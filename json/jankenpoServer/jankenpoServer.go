package main

import (
	"encoding/json"
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

	// cria um cofificador/decodificador Json
	jsonDecoder := json.NewDecoder(conn)
	jsonEncoder := json.NewEncoder(conn)
	var msgFromClient shared.Request

	fmt.Println("Servidor pronto para receber solicitações (TCP)...")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente e decodifica-as
		err = jsonDecoder.Decode(&msgFromClient)
		if err != nil {
			fmt.Println(err)
			return
		}

		// processa a solicitação
		r := shared.ProcessaSolicitacao(msgFromClient)

		// envia resposta ao cliente
		msgToClient := shared.Reply{r}
		err = jsonEncoder.Encode(msgToClient)
		if err != nil {
			fmt.Println(err)
			return
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
