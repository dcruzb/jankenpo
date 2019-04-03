package main

import (
	"flag"
	"fmt"
	rpcClient "jankenpo/impl/RPC/client"
	rmqClient "jankenpo/impl/rabbitMQ/client"
	jsonClient "jankenpo/impl/socketJson/client"
	tcpClient "jankenpo/impl/socketTCP/client"
	udpClient "jankenpo/impl/socketUDP/client"
	"jankenpo/shared"
	"sync"
	"time"
)

func main() {
	tcp := flag.Bool("tcp", shared.SOCKET_TCP, "Identifies if TCP client should start")
	udp := flag.Bool("udp", shared.SOCKET_UDP, "Identifies if UDP client should start")
	json := flag.Bool("json", shared.JSON, "Identifies if Json over TCP client should start")
	rpc := flag.Bool("rpc", shared.RPC, "Identifies if RPC client should start")
	rmq := flag.Bool("rmq", shared.RABBIT_MQ, "Identifies if RabbitMQ client should start")
	auto := flag.Bool("auto", shared.AUTO, "Identifies if the program should play in 'Auto' mode")
	flag.Parse()

	var wg sync.WaitGroup
	var elapsedTCP time.Duration
	var elapsedUDP time.Duration
	var elapsedJson time.Duration
	var elapsedRPC time.Duration
	var elapsedRMQ time.Duration

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

	fmt.Println("Tempo UDP:", elapsedUDP)
	fmt.Println("Tempo TCP:", elapsedTCP)
	fmt.Println("Tempo Json:", elapsedJson)
	fmt.Println("Tempo RPC:", elapsedRPC)
	fmt.Println("Tempo RabbitMQ:", elapsedRMQ)
}
