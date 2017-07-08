// ======================================================================================
// tcpClient.go  :  Creates a TCP connection to the TCP server on quic-proxy-client.go.
// Author        :  Donald Luc
// Date          :  July 7th, 2017
// ======================================================================================


/* Package */
package main


/* Imports */
import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"strings"
)


/* Constants */
const (
	IP   = "127.0.0.1"
	PORT = ":8686"
)


/* main */
func main() {
	log.Println("TCP Client '" + IP + PORT + "':")
	err := tcpClient()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("TCP Client is complete")
	os.Exit(0)
}


/* tcpClient */
func tcpClient() (error) {
	log.Println("Resolving remote server address...")
	addr, err := net.ResolveTCPAddr("tcp", IP+PORT)
	if err != nil { return errors.New("ResolveTCPAddr Failed: " + err.Error()) }

	log.Println("Dialing remote server address...")
	conn, err := net.DialTCP("tcp", nil, addr)
	defer conn.Close()
	if err != nil { return errors.New("Dial Failed: " + err.Error()) }

	log.Println("Writing to remote server...")
	log.Println("Please Type New Client Message:")
	reader := bufio.NewReader(os.Stdin)

	msg, _ := reader.ReadString('\n')
	msg = strings.Trim(msg, "\n ")

	_, err = conn.Write([]byte(msg))
	if err != nil { return errors.New("Write Failed: " + err.Error()) }

	log.Println("TCP Client successfully sent '" + msg + "'")
	return nil
}
