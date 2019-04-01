package main

import (
	"flag"
	jsonServer "jankenpo/impl/socketJson/server"
	tcpServer "jankenpo/impl/socketTCP/server"
	udpServer "jankenpo/impl/socketUDP/server"
	"sync"
)

func main() {
	tcp := flag.Bool("tcp", true, "Identifies if TCP server should start")
	udp := flag.Bool("udp", false, "Identifies if UDP server should start")
	json := flag.Bool("json", false, "Identifies if Json over TCP server should start")
	rpc := flag.Bool("rpc", false, "Identifies if RPC server should start")
	rmq := flag.Bool("rmq", false, "Identifies if RabbitMQ server should start")
	flag.Parse()

	var wg sync.WaitGroup

	if *tcp {
		wg.Add(1)
		go func() {
			tcpServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *udp {
		wg.Add(1)
		go func() {
			udpServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *json {
		wg.Add(1)
		go func() {
			jsonServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *rpc {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}

	if *rmq {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}

	wg.Wait()
}
