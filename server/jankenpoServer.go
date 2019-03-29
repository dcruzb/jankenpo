package main

import (
	"flag"
	"jankenpo/impl/socketJson"
	tcpServer "jankenpo/impl/socketTCP/server"
	"jankenpo/impl/socketUDP"
	"sync"
)

func main() {
	tcp := flag.Bool("tcp", true, "Identifies if TCP server should start")
	udp := flag.Bool("udp", true, "Identifies if UDP server should start")
	json := flag.Bool("json", true, "Identifies if Json over TCP server should start")
	rpc := flag.Bool("rpc", true, "Identifies if RPC server should start")
	rmq := flag.Bool("rmq", true, "Identifies if RabbitMQ server should start")
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
			socketUDP.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *json {
		wg.Add(1)
		go func() {
			socketJson.StartJankenpoServer()
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
