package main

import (
	"flag"
	"fmt"
	rpcClient "github.com/dcbCIn/jankenpo/impl/RPC/client"
	quicClient "github.com/dcbCIn/jankenpo/impl/quic/client"
	rmqClient "github.com/dcbCIn/jankenpo/impl/rabbitMQ/client"
	jsonClient "github.com/dcbCIn/jankenpo/impl/socketJson/client"
	tcpClient "github.com/dcbCIn/jankenpo/impl/socketTCP/client"
	udpClient "github.com/dcbCIn/jankenpo/impl/socketUDP/client"
	"github.com/dcbCIn/jankenpo/shared"
	"sync"
	"time"
)

func main() {
	quic := flag.Bool("quic", shared.QUIC, "Identifies if TCP client should start")
	tcp := flag.Bool("tcp", shared.SOCKET_TCP, "Identifies if TCP client should start")
	udp := flag.Bool("udp", shared.SOCKET_UDP, "Identifies if UDP client should start")
	json := flag.Bool("json", shared.JSON, "Identifies if Json over TCP client should start")
	rpc := flag.Bool("rpc", shared.RPC, "Identifies if RPC client should start")
	rmq := flag.Bool("rmq", shared.RABBIT_MQ, "Identifies if RabbitMQ client should start")
	auto := flag.Bool("auto", shared.AUTO, "Identifies if the program should play in 'Auto' mode")
	flag.Parse()

	var wg sync.WaitGroup
	var elapsedQuic time.Duration
	var elapsedTCP time.Duration
	var elapsedUDP time.Duration
	var elapsedJson time.Duration
	var elapsedRPC time.Duration
	var elapsedRMQ time.Duration

	if *quic {
		wg.Add(1)
		go func() {
			elapsedQuic = quicClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	if *tcp {
		wg.Add(1)
		go func() {
			elapsedTCP = tcpClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	if *udp {
		wg.Add(1)
		go func() {
			elapsedUDP = udpClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	if *json {
		wg.Add(1)
		go func() {
			elapsedJson = jsonClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	if *rpc {
		wg.Add(1)
		go func() {
			elapsedRPC = rpcClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	if *rmq {
		wg.Add(1)
		go func() {
			elapsedRMQ = rmqClient.PlayJanKenPo(*auto)
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Calls:", shared.SAMPLE_SIZE)
	fmt.Println("Wait:", shared.WAIT, "ms")
	fmt.Println("Tempo UDP:", elapsedUDP)
	fmt.Println("Tempo TCP:", elapsedTCP)
	fmt.Println("Tempo Quic:", elapsedQuic)
	fmt.Println("Tempo Json:", elapsedJson)
	fmt.Println("Tempo RPC:", elapsedRPC)
	fmt.Println("Tempo RabbitMQ:", elapsedRMQ)
}
