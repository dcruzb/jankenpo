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

func (st *SocketTCP) StartServer(ip, port string, useJson bool, initialConnections int) {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting TCP server. Details: ", err)
	}
	st.listener = ln
	st.useJson = useJson
	st.initialConnections = initialConnections
	st.clients = make([]Client, st.initialConnections)
}

func (st *SocketTCP) StopServer() {
	err := st.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (st *SocketTCP) ConnectToServer(ip, port string) {
	// connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		shared.PrintlnError(NAME, err)
	}

	st.serverConnection = conn
}

func (st *SocketTCP) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	conn, err := st.listener.Accept()
	if err != nil {
		shared.PrintlnError(NAME, "Error while waiting for connection", err)
	}

	cl = &st.clients[cliIdx]

	cl.connection = conn

	if st.useJson {
		// cria um cofificador/decodificador Json
		st.jsonDecoder = json.NewDecoder(conn)
		st.jsonEncoder = json.NewEncoder(conn)
	}

	return cl
}

func (st *SocketTCP) CloseConnection() {
	err := st.serverConnection.Close()
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

func (st *SocketTCP) Read() (message string) {
	if st.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from socket TCP. Details:", err)
		}
	}

	return message
}

func (st *SocketTCP) Write(message string) {
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

	_, err := st.serverConnection.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
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
