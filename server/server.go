package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
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
		req, err := ParseRequest(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		switch req.Command {
		case "addMember":
			if !cluster.MemberExists(req.Target) {
				addReq := Request{
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Target:  req.Target,
					Command: "addAck",
				}
				fmt.Printf("Attempting to reach `%s`.", req.Target)
				addReq.Send()
			} else {
				fmt.Printf("%s is already a member.", req.Target)
			}
			conn.Close()
		case "addAck":
			ackReq := Request{
				Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
				Target:  req.Source,
				Command: "addAccept",
			}
			fmt.Printf("%s established a connection with you.\n", req.Source)
			ackReq.Send()
		case "addAccept":
			member, _ := NewMember(req.Source)
			cluster = cluster.addMember(member)
			fmt.Printf("\n%s has joined the cluster.\n", member.Hostname)
			conn.Close()
		// case "sharedFiles":
		// 	if len(dataMap) == 3 {
		// 		member, err := cluster.FindMember(dataMap[1], dataMap[2])
		// 		if err != nil {
		// 			fmt.Printf("Member does not belong to cluster! %s", err)
		// 			return
		// 		}
		// 		member.Message(fmt.Sprintf("returnFileList %s %s", dataMap[1], dataMap[2]))
		// 	}
		// case "returnFileList":
		// 	var files []string
		// 	member, err := cluster.FindMember(dataMap[1], dataMap[2])
		// 	root := config.SharedDirectory
		// 	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 		files = append(files, path)
		// 		return nil
		// 	})
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	member.Message(fmt.Sprintf("listFiles %s %s", dataMap[1], dataMap[2]))
		//
		// case "listFiles":
		// 	fmt.Println(data)
		case "listNodes":
			fmt.Println(cluster.MemberHosts())
			conn.Close()
		case "ping":
			cluster.Ping()
			conn.Close()
		case "pong":
			fmt.Println("PONG")
		default:
			fmt.Printf("DAta: `%s`", req.String())
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
