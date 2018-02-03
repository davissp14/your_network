package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"your_network/configuration"
	"your_network/network"
)

var mynetwork network.Network

func Start(config configuration.Configuration) {
	mynetwork.Config = config
	fmt.Println("Starting Server on Port: " + config.Port)
	var listener, err = ListenerInterface(config.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initiate Network
	//mynetwork.Init()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		data, _ := reader.ReadString('\n')
		if data == "" {
			break
		}
		req, err := network.ParseRequest(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		mynetwork.HandleRequest(network.Connection{Conn: conn}, req)
		if req.Client {
			conn.Close()
		}
	}
}

func ListenerInterface(port string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+port)
	if err != nil {
		return nil, errors.New("Failed to resolve tcp addr")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, errors.New("Failed to listen to tcp")
	}
	return listener, nil
}
