package socketTcpSsl

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"github.com/dcbCIn/jankenpo/shared"
	"net"
	"os"
)

const NAME = "jankenpo/socketTcpSsl"

type Client struct {
	connection net.Conn
}

type SocketTcpSsl struct {
	ip                 string
	port               string
	useJson            bool
	listener           net.Listener
	serverConnection   *tls.Conn
	initialConnections int
	clients            []Client

	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

func (st *SocketTcpSsl) StartServer(ip, port string, useJson bool, initialConnections int) {
	ln, err := tls.Listen("tcp4", ip+":"+port, shared.GetServerTLSConfig())
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting TcpSsl server. Details: ", err)
	}
	st.listener = ln
	st.useJson = useJson
	st.initialConnections = initialConnections
	st.clients = make([]Client, st.initialConnections)
}

func (st *SocketTcpSsl) StopServer() {
	err := st.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (st *SocketTcpSsl) ConnectToServer(ip, port string) {
	conn, err := tls.Dial("tcp4", ip+":"+port, shared.GetClientTLSConfig())
	if err != nil {
		shared.PrintlnError(NAME, err)
	}

	st.serverConnection = conn
}

func (st *SocketTcpSsl) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
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

func (st *SocketTcpSsl) CloseConnection() {
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

func (st *SocketTcpSsl) Read() (message string) {
	if st.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from socket TcpSsl. Details:", err)
		}
	}

	return message
}

func (st *SocketTcpSsl) Write(message string) {
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
		shared.PrintlnError(NAME, "Error while writing message to socket TcpSsl. Details:", err)
		os.Exit(1)
	}
}

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		shared.PrintlnError(NAME, "Error while reading message from socket TcpSsl. Details:", err)
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
		shared.PrintlnError(NAME, "Error while writing message to socket TcpSsl. Details:", err)
		os.Exit(1)
	}
}
