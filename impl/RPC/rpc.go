package rpc

import (
	"bufio"
	"encoding/json"
	"jankenpo/shared"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const NAME = "jankenpo/rpc"

type Client struct {
	connection net.Conn
}

type RPC struct {
	ip                 string
	port               string
	useJson            bool
	listener           net.Listener
	rpcClient          rpc.Client
	initialConnections int
	clients            []Client

	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func (this *RPC) StartServer(ip, port string, useJson bool, initialConnections int) {
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

	err = http.Serve(ln, nil) // TODO move to ConnectToServer
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting RPC server. Details: ", err)
	}

	this.listener = ln
	this.useJson = useJson
	this.initialConnections = initialConnections
	this.clients = make([]Client, this.initialConnections)
}

func (this *RPC) StopServer() {
	err := this.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (this *RPC) ConnectToServer(ip, port string) {
	// connect to server
	rpcClient, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	/*conn, err := net.Dial("tcp", ip +":"+ port)
	if err != nil {
		shared.PrintlnError(NAME, err)
	}*/

	this.rpcClient = *rpcClient
}

func (this *RPC) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	conn, err := this.listener.Accept()
	if err != nil {
		shared.PrintlnError(NAME, "Error while waiting for connection", err)
	}

	cl = &this.clients[cliIdx]

	cl.connection = conn

	if this.useJson {
		// cria um cofificador/decodificador Json
		this.jsonDecoder = json.NewDecoder(conn)
		this.jsonEncoder = json.NewEncoder(conn)
	}

	return cl
}

func (this *RPC) CloseConnection() {
	err := this.rpcClient.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
}

func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
}

func (this *RPC) Call(serviceMethod string, request shared.Request) (reply *shared.Reply) {

	this.rpcClient.Call(serviceMethod, request, &reply)

	return reply
}

/*func (rpc *RPC) Read() (message string) {
	if rpc.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		message, err = bufio.NewReader(rpc.serverConnection).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from socket TCP. Details:", err)
		}
	}

	return message
}*/

func (this *RPC) Write(message string) {
	// envia resporpca

	// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
	//_, err := fmt.Fprintf(conn, msgToServer+"\n")
	//_, err := conn.Write([]byte( msgToServer + "\n"))
	/*reader := bufio.NewWriter(conn)
	_, err := reader.WriteString( msgToServer + "\n")
	reader.Flush()*/
	/*reader := bufio.NewWriter(conn)
	_, err := io.WriteString(reader, msgToServer + "\n")
	reader.Flush()*/
	//_, err := io.WriteString(conn, msgToServer+"\n")

	/*_, err := rpc.serverConnection.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}*/
}

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		shared.PrintlnError(NAME, "Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (cl *Client) Write(message string) {
	// envia resposta

	// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
	//_, err := fmt.Fprintf(conn, msgToServer+"\n")
	//_, err := conn.Write([]byte( msgToServer + "\n"))
	/*reader := bufio.NewWriter(conn)
	_, err := reader.WriteString( msgToServer + "\n")
	reader.Flush()*/
	/*reader := bufio.NewWriter(conn)
	_, err := io.WriteString(reader, msgToServer + "\n")
	reader.Flush()*/
	//_, err := io.WriteString(conn, msgToServer+"\n")

	_, err := cl.connection.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}
