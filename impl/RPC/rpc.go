package rpc

import (
	"github.com/dcbCIn/jankenpo/shared"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

const NAME = "jankenpo/rpc"

type RPC struct {
	ip        string
	port      string
	listener  net.Listener
	rpcClient rpc.Client
}

func (this *RPC) StartServer(ip, port string) {
	request := new(shared.Request)
	// Publish the receivers methods
	err := rpc.Register(request)
	if err != nil {
		shared.PrintlnError(NAME, "Error while registering RPC receivers methods. Details: ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()

	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting RPC server. Details: ", err)
	}

	this.listener = ln
}

func (this *RPC) StopServer() {
	err := this.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (this *RPC) WaitForConnection(cliIdx int) {
	err := http.Serve(this.listener, nil)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting RPC server. Details: ", err)
	}
}

func (this *RPC) ConnectToServer(ip, port string) {
	// connect to server
	rpcClient, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	this.rpcClient = *rpcClient
}

func (this *RPC) CloseConnection() {
	err := this.rpcClient.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
}

func (this *RPC) Call(serviceMethod string, request shared.Request) (reply *shared.Reply) {
	err := this.rpcClient.Call(serviceMethod, request, &reply)
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
	return reply
}
