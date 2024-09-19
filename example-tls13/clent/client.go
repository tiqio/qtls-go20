package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	port := flag.String("port", "8360", "port to connect")
	certFile := flag.String("certfile", "../../testdata/example-cert.pem", "trusted CA certificate")
	//keyLogFilePath := flag.String("log", "./client.log", "tls key log path")
	flag.Parse()
	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	//keyLogFile, _ := os.OpenFile(*keyLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	//fmt.Fprintf(keyLogFile, "# SSL/TLS secrets log file, generated by go\n")

	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
		//KeyLogWriter:       keyLogFile,
	}
	conn, err := tls.Dial("tcp4", "localhost:"+*port, config)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(conn, "Hello simple secure Server\n")
	if err != nil {
		log.Fatal("client write error:", err)
	}
	if err = conn.CloseWrite(); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Println("client read:", string(buf[:n]))
	conn.Close()
}
