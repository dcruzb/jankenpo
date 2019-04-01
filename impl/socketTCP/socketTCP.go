package socketTCP

import (
	"bufio"
	"encoding/json"
	"jankenpo/shared"
	"net"
	"os"
)

const NAME = "jankenpo/socketTCP"

type Client struct {
	connection net.Conn
}

type SocketTCP struct {
	ip         string
	port       string
	useJson    bool
	listener   net.Listener
	connection net.Conn
	clients    []Client

	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func (st *SocketTCP) StartServer(ip, port string, useJson bool) {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting TCP server. Details: ", err)
	}
	st.listener = ln
	st.useJson = useJson
}

func (st *SocketTCP) StopServer() {
	err := st.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (st *SocketTCP) WaitForConnection() {
	// aceita conexões na porta
	conn, err := st.listener.Accept()
	if err != nil {
		shared.PrintlnError(NAME, "Error while waiting for connection", err)
	}
	st.connection = conn

	if st.useJson {
		// cria um cofificador/decodificador Json
		st.jsonDecoder = json.NewDecoder(conn)
		st.jsonEncoder = json.NewEncoder(conn)
	}
}

func (st *SocketTCP) CloseConnection() {
	err := st.connection.Close()
	if err != nil {
		shared.PrintlnError(NAME, err)
	}
}

func (st *SocketTCP) Read() (message string) {
	if st.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		message, err = bufio.NewReader(st.connection).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from socket TCP. Details:", err)
		}
	}

	return message
}

func (st *SocketTCP) Write(message string) {
	// envia resposta ao cliente
	_, err := st.connection.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}
