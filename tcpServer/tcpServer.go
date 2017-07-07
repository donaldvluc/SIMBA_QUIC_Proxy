// ======================================================================================
// tcpServer.go  :  Creates a TCP socket for the TCP client from quicServer.go
// Author        :  Donald Luc
// Date          :  July 6th, 2017
// ======================================================================================


// Package       :
package main


// Imports       :
import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
)


// Constants     :
const(
	PORT = ":8686"
)


// Main          :
func main() {
	log.Println("TCP Server '" + PORT + "'")
	err := tcpServer()

	if err != nil {
		log.Println("TCP server error: ", err)
		os.Exit(1)
	}

	log.Println("TCP server is complete")
	os.Exit(0)
}


// tcpServer     :
func tcpServer() (error) {
	log.Println("Resolving server address...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", PORT)
	if err != nil { return errors.New("ResolveTCPAddr Failed: " + err.Error()) }

	log.Println("Listening for incoming connections...")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil { return errors.New("ListenTCP Failed: " + err.Error()) }

	log.Println("Accepting all incoming connections...")

	conn, err := tcpListener.AcceptTCP()
	defer conn.Close()
	if err != nil { return errors.New("AcceptTCP Failed: " + err.Error()) }

	var buf bytes.Buffer
	io.Copy(&buf, conn)

	log.Println("TCP server successfully received '" + buf.String() + "'")

	return nil
}
