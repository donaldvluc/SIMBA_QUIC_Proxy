// ======================================================================================
// tcpServer.go  :  Creates a TCP socket for the TCP client from quicServer.go.
// Author        :  Donald Luc
// Date          :  July 7th, 2017
// ======================================================================================


/* Package */
package main


/* Imports */
import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
)


/* Constants */
const(
	PORT = ":8688"
)


/* main */
func main() {
	log.Println("TCP Server listening on port '" + PORT + "'...")
	err := tcpServer()

	if err != nil {
		log.Println("TCP Server Error -->", err)
		os.Exit(1)
	}

	log.Println("TCP Server is complete")
	os.Exit(0)
}


/* tcpServer */
func tcpServer() (error) {
	// Resolve the given address as TCP.
	log.Println("Resolving TCP Server address...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", PORT)
	if err != nil { return errors.New("ResolveTCPAddr Failed --> " + err.Error()) }

	// Listen for TCP streams.
	log.Println("Listening for incoming connections...")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil { return errors.New("ListenTCP Failed --> " + err.Error()) }

	// Accept TCP data.
	log.Println("Accepting all incoming connections...")
	conn, err := tcpListener.AcceptTCP()
	defer conn.Close()
	if err != nil { return errors.New("AcceptTCP Failed --> " + err.Error()) }

	log.Println("Handling incoming connections...")
	return readHandler(conn)
}


/* readHandler */
func readHandler(conn net.Conn) (error) {
	reader := bufio.NewReader(conn)

	// Read messages until EOF or error is reached.
	log.Println("Reading incoming messages...")
	msg, err := reader.ReadString('\n')

	// Handle EOF and errors.
	switch{

	case err == io.EOF:
		return errors.New("Reached EOF --> Closing Connection...")

	case err != nil:
		return errors.New("ReadString Error --> " + err.Error())

	}

	// Trim endline since ReadString does not remove it.
	trimMsg := strings.Trim(msg, "\n")
	log.Println("TCP Server successfully received:\n" + trimMsg)
	return nil
}