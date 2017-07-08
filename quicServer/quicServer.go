// ======================================================================================
// quicServer.go  :  Creates a QUIC server to accept client data to then forward to a TCP
//                   server at tcpServer.go
// Author         :  Donald Luc
// Date           :  July 7th, 2017
// ======================================================================================


/* Package */
package main


/* Imports */
import (
//	"bufio"
	"errors"
	"log"
	"net"
	"os"
//	"strings"


	"io"
	"math/big"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	quic "github.com/lucas-clemente/quic-go"
)


/* Constants */
const (
	QUIC_ADDR = "127.0.0.1:8688"
	TCP_ADDR  = "127.0.0.1:8688"
)


/* main */
func main() {
	log.Println("Setting up network proxy with QUIC Server and TCP Client...")

	// Phase 1 : QUIC Server
	log.Println("PHASE 1 : QUIC SERVER ===========================")

	msg, err := quicServer()
	if err != nil {
		log.Println("QUIC Server Error -->", err)
		os.Exit(1)
	}
	log.Println("QUIC Server is complete")

	// Phase 2 : TCP Client
	log.Println("PHASE 2 : TCP CLIENT  ===========================")

	log.Println("TCP Client '" + TCP_ADDR + "':")
	err = tcpClient(msg)
	if err != nil {
		log.Println("TCP Client Error -->", err)
		os.Exit(1)
	}
	log.Println("TCP Client is complete")

	log.Println("Closing QUIC Proxy...")
	os.Exit(0)
}


/* quicServer */
func quicServer() (string, error) {
	log.Println("Setting up QUIC Server configurations...")
	config := &quic.Config{
		TLSConfig: generateTLSConfig(),
	}

	log.Println("Listening for incoming QUIC connections...")
	listener, err := quic.ListenAddr(QUIC_ADDR, config)
	if err != nil { return "", errors.New("QUIC ListenAddr Failed --> " + err.Error()) }

	log.Println("Accepting all incoming QUIC connections...")
	session, err := listener.Accept()
	if err != nil { return "", errors.New("QUIC Accept Failed --> " + err.Error()) }

	log.Println("Accepting QUIC stream...")
	stream, err := session.AcceptStream()
	if err != nil { return "", errors.New("QUIC AcceptStream Failed --> " + err.Error()) }

	log.Println("Handling QUIC stream...")
	msg, err := handleStream(stream)
	if err != nil { return "", errors.New("QUIC handleStream Failed --> " + err.Error()) }
	
	log.Println("QUIC Server successfully received '" + msg + "'")
	return msg, nil
}


/* handleStream */
func handleStream(stream quic.Stream) (string, error) {
	buf := make([]byte, 10)
	_, err := io.ReadFull(stream, buf)
	if err != nil { return "", err }

	msg := string(buf)
	return msg, nil
}


// generateTLSConfig
func generateTLSConfig() (*tls.Config) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}


/* tcpClient */
func tcpClient(msg string) (error) {
	log.Println("Resolving remote server address...")
	addr, err := net.ResolveTCPAddr("tcp", TCP_ADDR)
	if err != nil { return errors.New("ResolveTCPAddr Failed: " + err.Error()) }

	log.Println("Dialing remote server address...")
	conn, err := net.DialTCP("tcp", nil, addr)
	defer conn.Close()
	if err != nil { return errors.New("Dial Failed: " + err.Error()) }

	log.Println("Forwarding message to remote server...")
	_, err = conn.Write([]byte(msg))
	if err != nil { return errors.New("Write Failed: " + err.Error()) }

	log.Println("TCP Client successfully sent '" + msg + "'")
	return nil
}