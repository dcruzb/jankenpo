package main

import (
	"flag"
	rpcServer "github.com/dcbCIn/jankenpo/impl/RPC/server"
	quicServer "github.com/dcbCIn/jankenpo/impl/quic/server"
	rmqServer "github.com/dcbCIn/jankenpo/impl/rabbitMQ/server"
	jsonServer "github.com/dcbCIn/jankenpo/impl/socketJson/server"
	tcpServer "github.com/dcbCIn/jankenpo/impl/socketTCP/server"
	tcpSslServer "github.com/dcbCIn/jankenpo/impl/socketTcpSsl/server"
	udpServer "github.com/dcbCIn/jankenpo/impl/socketUDP/server"
	"github.com/dcbCIn/jankenpo/shared"
	"sync"
)

func main() {
	quic := flag.Bool("quic", shared.QUIC, "Identifies if Quic server should start")
	tcp := flag.Bool("tcp", shared.SOCKET_TCP, "Identifies if TCP server should start")
	tcpSsl := flag.Bool("tcpSsl", shared.SOCKET_TCP_SSL, "Identifies if TCP SSL server should start")
	udp := flag.Bool("udp", shared.SOCKET_UDP, "Identifies if UDP server should start")
	json := flag.Bool("json", shared.JSON, "Identifies if Json over TCP server should start")
	rpc := flag.Bool("rpc", shared.RPC, "Identifies if RPC server should start")
	rmq := flag.Bool("rmq", shared.RABBIT_MQ, "Identifies if RabbitMQ server should start")
	flag.Parse()

	var wg sync.WaitGroup

	if *quic {
		wg.Add(1)
		go func() {
			quicServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *tcp {
		wg.Add(1)
		go func() {
			tcpServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *tcpSsl {
		wg.Add(1)
		go func() {
			tcpSslServer.StartJankenpoServer()
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
			rpcServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	if *rmq {
		wg.Add(1)
		go func() {
			rmqServer.StartJankenpoServer()
			wg.Done()
		}()
	}

	wg.Wait()
}
