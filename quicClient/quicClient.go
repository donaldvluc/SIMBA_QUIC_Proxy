// ======================================================================================
// quicClient.go  :  Creates a TCP server to accept client data to then forward to a QUIC
//                   server at quicServer.go
// Author         :  Donald Luc
// Date           :  July 7th, 2017
// ======================================================================================


/* Package */
package main


/* Imports */
import (
	"bufio"
	"crypto/tls"  // QUIC Client
	"errors"
	"io"
	"log"
	"net"         // TCP Server 
	"os"
	"strings"

	quic "github.com/lucas-clemente/quic-go"
)


/* Constants */
const(
	TCP_PORT  = ":8686"
	QUIC_ADDR = "127.0.0.1:8688"
)


/* main */
func main() {
	log.Println("Setting up network proxy with TCP Server and QUIC Client...")

	// Phase 1 - TCP Server Setup
	log.Println("PHASE 1 : TCP SERVER  =====================================")

	log.Println("TCP Server listening on address '" + TCP_PORT + "'...")
	msg, err := tcpServer()
	if err != nil {
		log.Println("TCP Server Error -->", err)
		os.Exit(1)
	}
	log.Println("TCP Server completed...")

	// Phase 2 - QUIC Client Setup
	log.Println("PHASE 2 : QUIC CLIENT =====================================")

	log.Println("QUIC Client writing to address '" + QUIC_ADDR + "'...")
	err = quicClient(msg)
	if err != nil {
		log.Println("QUIC Client Error -->", err)
		os.Exit(1)
	}
	log.Println("QUIC Client completed...")

	log.Println("Closing QUIC Proxy...")
	os.Exit(0)
}


/* tcpServer */
func tcpServer() (string, error) {
	log.Println("Resolving TCP Server address...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", TCP_PORT)
	if err != nil { return "", errors.New("ResolveTCPAddr Failed --> " + err.Error()) }

	log.Println("Listening for incoming TCP connections...")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil { return "", errors.New("ListenTCP Failed --> " + err.Error()) }

	log.Println("Accepting all incoming TCP connections...")
	conn, err := tcpListener.AcceptTCP()
	defer conn.Close()
	if err != nil { return "", errors.New("AcceptTCP Failed --> " + err.Error()) }

	return readHandler(conn)
}


/* readHandler */
func readHandler(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	log.Println("Reading incoming TCP messages...")
	msg, err := reader.ReadString('\n')

	switch{
	
	case err == io.EOF:
		return "", errors.New("Reached EOF --> Closing Connection...")
	
	case err != nil:
		return "", errors.New("ReadString Error --> " + err.Error())

	}

	trimMsg := strings.Trim(msg, "\n")
	log.Println("TCP Server successfully received '" + trimMsg + "'")
	return msg, nil
}


/* quicClient */
func quicClient(msg string) (error) {
	log.Println("Setting up QUIC Client configurations...")
	config := &quic.Config{
		TLSConfig: &tls.Config{ InsecureSkipVerify: true },
	}

	log.Println("Dialing QUIC server...")
	session, err := quic.DialAddr(QUIC_ADDR, config)
	if err != nil { return errors.New("QUIC DialAddr Failed --> " + err.Error()) }

	log.Println("Synchronizing stream between QUIC Client and QUIC Server...")
	stream, err := session.OpenStreamSync()
	if err != nil { return errors.New("QUIC OpenStreamSync Failed --> " + err.Error()) }

	return writeHandler(stream, msg)
}


/* writeHandler */
func writeHandler(stream quic.Stream, msg string) (error) {
	log.Println("Writing to QUIC server...")
	_, err := stream.Write([]byte(msg))
	if err != nil { return errors.New("QUIC Write Failed --> " + err.Error()) }

	trimMsg := strings.Trim(msg, "\n")
	log.Println("QUIC Client successfully sent '" + trimMsg + "'")
	return nil
}