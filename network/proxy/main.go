package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"net"
)

const (
	MODULE  = "Chilopod Network Module"
	VERSION = "-1.0-release-201812"
	SERVICE = "127.0.1.1:8080"
	AUTHOR  = "@RiboseYim"
)

func main() {

	local_ip := flag.String("lip", "", "local addr")
	local_port := flag.String("lport", "", "local port")
	remote_ip := flag.String("rip", "", "remote addr")
	remote_port := flag.String("rport", "", "remote addr")

	flag.Parse()

	local_addr := *local_ip + ":" + *local_port
	remote_addr := *remote_ip + ":" + *remote_port

	log.Printf("Welcome to [ Chilopod System %s ] \n Author:%s \n\n", VERSION, AUTHOR)
	log.Println("local_ip:", *local_ip)
	log.Println("local_port:", *local_port)
	log.Println("remote_ip:", *remote_ip)
	log.Println("remote_port:", *remote_port)
	log.Println("local_addr:", local_addr)
	log.Println("remote_addr:", remote_addr)

	addr, err := net.ResolveTCPAddr("tcp", local_addr)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

	for i := 0; i < 5; i++ {
		go handleConn(*remote_ip, *remote_port, pending, complete)
	}
	go closeConn(complete)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("-----AcceptTCP err:%s", err)
			panic(err)
		}
		pending <- conn
	}
}

func handleConn(remote_ip string, remote_port string, in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	remote_addr := remote_ip + ":" + remote_port
	log.Println("-----handleConn:" + remote_addr)
	for conn := range in {
		proxyConn(remote_addr, conn)
		out <- conn
	}
}

func proxyConn(remote_addr string, conn *net.TCPConn) {
	rAddr, err := net.ResolveTCPAddr("tcp", remote_addr)
	if err != nil {
		log.Println("-----proxyConn rAddr ResolveTCPAddr err:%s", err)
		panic(err)
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		log.Println("-----proxyConn rAddr DialTCP err:%s", err)
		panic(err)
	}
	defer rConn.Close()

	buf := &bytes.Buffer{}
	for {
		data := make([]byte, 256)
		n, err := conn.Read(data)
		if err != nil {
			log.Println("-----conn.Read(data) err:%s", err)
			panic(err)
		}
		buf.Write(data[:n])
		if data[0] == '\r' && data[1] == '\n' {
			break
		}
	}

	if _, err := rConn.Write(buf.Bytes()); err != nil {
		panic(err)
	}

	log.Printf("sent:\n%v", hex.Dump(buf.Bytes()))

	data := make([]byte, 1024)
	n, err := rConn.Read(data)
	if err != nil {
		if err != io.EOF {
			log.Println("-----err != io.EOF err:%s", err)
			panic(err)
		} else {
			log.Printf("received err: %v", err)
		}
	}
	log.Printf("received:\n%v", hex.Dump(data[:n]))
}

func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}
