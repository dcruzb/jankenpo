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

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ProcessaSolicitacao(request shared.Request) int {
	possibilities := []string{"P", "A", "T"}

	if !inArray(request.Player1, possibilities) {
		return -1
	}
	if !inArray(request.Player2, possibilities) {
		return -1
	}
	if request.Player1 == request.Player2 {
		return 0
	}
	switch request.Player1 {
	case "P":
		if request.Player2 == "A" {
			return 2
		} else {
			return 1
		}
	case "A":
		if request.Player2 == "P" {
			return 1
		} else {
			return 2
		}
	case "T":
		if request.Player2 == "P" {
			return 2
		} else {
			return 1
		}
	default:
		return -1
	}
}

func StartJankenpoServer() {
	fmt.Println("Initializing server")

	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(shared.SERVER_PORT))
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		os.Exit(1)
	}

	// escuta na porta tcp 4600
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
	message := make([]byte, 8)

	fmt.Println("Servidor pronto para receber solicitações (UDP)...")
	for idx := 0; idx < shared.SAMPLE_SIZE; idx++ {

		// recebe solicitações do cliente
		n, addr, err := conn.ReadFromUDP(message)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Message received: ", message)
		_, err = fmt.Sscanf(string(message[:n]), "%s %s", &msgFromClient.Player1, &msgFromClient.Player2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// processa a solicitação
		r := ProcessaSolicitacao(msgFromClient)

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
