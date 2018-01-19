package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"your_network/configuration"
)

var cluster Cluster
var config configuration.Configuration

func Start(conf configuration.Configuration) {
	config = conf
	fmt.Println("Starting Server on Port: " + config.Port)
	var listener, err = ListenerInterface(config.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
		dataMap := strings.Split(data, " ")
		switch strings.TrimSpace(dataMap[0]) {
		case "addMember":
			if len(dataMap) == 3 {
				member, err := NewMember(dataMap[1], dataMap[2])
				if err != nil {
					fmt.Printf("Error occurred while adding member %s", err)
					return
				}
				if !cluster.MemberExists(member) {
					fmt.Printf("Attempting to reach `%s`.", member.Hostname)
					member.Message(fmt.Sprintf("addAck %s %s", config.Hostname, config.Port))
					member.Conn.Close()
				} else {
					fmt.Printf("%s is already a member.", member.Hostname)
				}
			} else {
				fmt.Println("Usage: `addMember localhost <port>`")
			}
			conn.Close()
		case "addAck":
			if len(dataMap) == 3 {
				member, err := NewMember(dataMap[1], dataMap[2])
				if err != nil {
					fmt.Printf("Error occurred while adding member %s", err)
					return
				}
				fmt.Printf("%s established a connection with you. Success!", member.Hostname)
				member.Message(fmt.Sprintf("addAccept %s %s", config.Hostname, config.Port))
				member.Conn.Close()
			} else {
				fmt.Println("I'm down to pair, but you need to send me the right shit.")
			}
			conn.Close()
		case "addAccept":
			if len(dataMap) == 3 {
				member, err := NewMember(dataMap[1], dataMap[2])
				if err != nil {
					fmt.Printf("Error occurred while adding member %s", err)
					return
				}
				cluster = cluster.addMember(member)
				fmt.Printf("\n%s has joined the cluster.", member.Hostname)
			} else {
				fmt.Println("Accept message is invalid!")
			}
			conn.Close()
		case "message":
			fmt.Println(len(cluster.Members))
			cluster.Communicate("received\n")
			conn.Close()
		case "received":
			fmt.Println("Wave")
		default:
			fmt.Println(data)
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
