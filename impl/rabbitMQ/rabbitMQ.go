package rabbitMQ

import (
	"bufio"
	"encoding/json"
	"jankenpo/shared"
	"net"
	"os"
)

const NAME = "jankenpo/rabbitMQ"

type Client struct {
	connection net.Conn
}

type RabbitMQ struct {
	ip                 string
	port               string
	useJson            bool
	listener           net.Listener
	serverConnection   net.Conn
	initialConnections int
	clients            []Client

	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func (rMQ *RabbitMQ) StartServer(ip, port string, useJson bool, initialConnections int) {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting rabbitMQ server. Details: ", err)
	}
	rMQ.listener = ln
	rMQ.useJson = useJson
	rMQ.initialConnections = initialConnections
	rMQ.clients = make([]Client, rMQ.initialConnections)
}

func (rMQ *RabbitMQ) StopServer() {
	err := rMQ.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (rMQ *RabbitMQ) ConnectToServer(ip, port string) {
	// connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, err)
	}

	rMQ.serverConnection = conn
}

func (rMQ *RabbitMQ) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	conn, err := rMQ.listener.Accept()
	if err != nil {
		shared.PrintlnError(NAME, "Error while waiting for connection", err)
	}

	cl = &rMQ.clients[cliIdx]

	cl.connection = conn

	if rMQ.useJson {
		// cria um cofificador/decodificador Json
		rMQ.jsonDecoder = json.NewDecoder(conn)
		rMQ.jsonEncoder = json.NewEncoder(conn)
	}

	return cl
}

func (rMQ *RabbitMQ) CloseConnection() {
	err := rMQ.serverConnection.Close()
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

func (rMQ *RabbitMQ) Read() (message string) {
	if rMQ.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		message, err = bufio.NewReader(rMQ.serverConnection).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from rabbitMQ. Details:", err)
		}
	}

	return message
}

func (rMQ *RabbitMQ) Write(message string) {
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

	_, err := rMQ.serverConnection.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to rabbitMQ. Details:", err)
		os.Exit(1)
	}
}

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		shared.PrintlnError(NAME, "Error while reading message from rabbitMQ. Details:", err)
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
		shared.PrintlnError(NAME, "Error while writing message to rabbitMQ. Details:", err)
		os.Exit(1)
	}
}
