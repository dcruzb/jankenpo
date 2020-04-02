package quic

import (
	"bufio"
	"context"
	"crypto/tls"
	//"encoding/json"
	"github.com/dcbCIn/jankenpo/shared"
	"github.com/lucas-clemente/quic-go"
	"os"
)

const NAME = "jankenpo/quic"

type Client struct {
	session 	quic.Session //net.Conn
	stream 		quic.Stream
}

type Quic struct {
	ip                 	string
	port               	string
	useJson       		bool
	listener           	quic.Listener
	serverSession	   	quic.Session
	initialConnections	int
	clients        		[]Client
	stream				quic.Stream
	//jsonEncoder *json.Encoder
	//jsonDecoder *json.Decoder
}

func (st *Quic) StartServer(ip, port string, useJson bool, initialConnections int) {
	//ln, err := net.Listen("tcp", ip+":"+port)
	quicConfig := quic.Config{ KeepAlive:true}
	ln, err := quic.ListenAddr(ip + ":" + port, shared.GenerateTLSConfig(), &quicConfig)
	if err != nil {
		shared.PrintlnError(NAME, "Error while starting TCP server. Details: ", err)
	}
	st.listener = ln
	st.useJson = useJson
	st.initialConnections = initialConnections
	st.clients = make([]Client, st.initialConnections)
}

func (st *Quic) StopServer() {
	err := st.listener.Close()
	if err != nil {
		shared.PrintlnError(NAME, "Error while stoping server. Details:", err)
	}
}

func (st *Quic) ConnectToServer(ip, port string) {
	// connect to server
	//conn, err := net.Dial("tcp", ip+":"+port)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"exemplo"},
	}
	session, err := quic.DialAddr(ip+":"+port, tlsConf, nil)
	if err != nil {
		shared.PrintlnError(NAME, err)
	}

	st.serverSession = session

	if st.stream == nil {
		stream, err := st.serverSession.OpenStreamSync(context.Background())
		st.stream = stream
		if err != nil {
			shared.PrintlnError(NAME, "Error while writing message to quic. Details:", err)
			os.Exit(1)
		}
	}
}

func (st *Quic) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	session, err := st.listener.Accept(context.Background())
	if err != nil {
		shared.PrintlnError(NAME, "Error while waiting for connection", err)
	}

	cl = &st.clients[cliIdx]

	cl.session = session

	if cl.stream == nil {
		stream, err := cl.session.AcceptStream(context.Background())
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading stream from quic. Details:", err)
		}
		cl.stream = stream
	}


	//if st.useJson {
	//	// cria um cofificador/decodificador Json
	//	st.jsonDecoder = json.NewDecoder(session)
	//	st.jsonEncoder = json.NewEncoder(session)
	//}

	return cl
}

//func (st *SocketTCP) CloseConnection() {
//	err := st.serverSession.Close()
//	if err != nil {
//		shared.PrintlnError(NAME, err)
//	}
//}
//
//func (cl *Client) CloseConnection() {
//	err := cl.session.Close()
//	if err != nil {
//		shared.PrintlnError(NAME, err)
//	}
//}

func (st *Quic) Read() (message string) {
	if st.useJson {

	} else {
		var err error
		// recebe solicitações do cliente
		//message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
		//if st.stream == nil {
		//	stream, err := st.serverSession.AcceptStream(context.Background())
		//	if err != nil {
		//		shared.PrintlnError(NAME, "Error while reading stream from quic. Details:", err)
		//	}
		//	st.stream = stream
		//}
		message, err = bufio.NewReader(st.stream).ReadString('\n')
		if err != nil {
			shared.PrintlnError(NAME, "Error while reading message from quic stream. Details:", err)
		}
	}

	return message
}

func (st *Quic) Write(message string) {
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

	//_, err := st.serverSession.Write([]byte(message + "\n"))

	//if st.stream == nil {
	//	stream, err := st.serverSession.OpenStreamSync(context.Background())
	//	st.stream = stream
	//	if err != nil {
	//		shared.PrintlnError(NAME, "Error while writing message to quic. Details:", err)
	//		os.Exit(1)
	//	}
	//}

	_, err := st.stream.Write([]byte(message + "\n"))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to quic. Details:", err)
		os.Exit(1)
	}
}

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	//message, err = bufio.NewReader(cl.connection).ReadString('\n')
	//if cl.stream == nil {
	//	stream, err := cl.session.AcceptStream(context.Background())
	//	if err != nil {
	//		shared.PrintlnError(NAME, "Error while reading stream from quic. Details:", err)
	//	}
	//	cl.stream = stream
	//}

	message, err = bufio.NewReader(cl.stream).ReadString('\n')

	if err != nil {
		shared.PrintlnError(NAME, "Error at client while reading message from the server stream. Details:", err)
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

	//_, err := cl.connection.Write([]byte(message + "\n"))

	if cl.stream == nil {
		stream, err := cl.session.OpenStreamSync(context.Background())
		cl.stream = stream
		if err != nil {
			shared.PrintlnError(NAME, "Error while writing message to quic. Details:", err)
			os.Exit(1)
		}
	}

	_, err := cl.stream.Write([]byte(message))
	if err != nil {
		shared.PrintlnError(NAME, "Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}

