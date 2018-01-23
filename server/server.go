package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"your_network/configuration"
	"your_network/network"
)

var mynetwork network.Network
var config configuration.Configuration

func Start(conf configuration.Configuration) {
	config = conf
	fmt.Println("Starting Server on Port: " + config.Port)
	var listener, err = ListenerInterface(config.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initiate Network
	// mynetwork.Init(config)

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
		// fmt.Println(req)
		switch req.Command {
		// WIP possible solution to initiating an existing network on boot.
		case "init_connection":
			mynetwork = mynetwork.AddNode(req.Target)
			conn.Close()
		// Will work to pair with a remote network.
		case "membershipRequest":
			switch req.State {
			case "new":
				memReq := network.Request{
					Target:  req.Target,
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Command: "membershipRequest",
					State:   "request",
				}
				memReq.Send()
				conn.Close()
			case "request":
				memReq := network.Request{
					Target:  req.Source,
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Command: "membershipRequest",
					State:   "response",
					Success: true,
				}
				if !mynetwork.NodeExists(req.Source) {
					mynetwork = mynetwork.AddNode(req.Source)
				}
				node, err := mynetwork.FindNode(req.Source)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("\nYou and `%s` have been paired!", req.Source)
				node.SendRequest(memReq)
				conn.Close()

			case "response":
				if req.Success {
					mynetwork = mynetwork.AddNode(req.Source)
					fmt.Printf("\nYou and `%s` have been paired!", req.Source)
				} else {
					fmt.Println("Failed to add Node! :(")
				}
			}
		// Lists all of the nodes in the network.
		case "listNodes":
			for _, node := range mynetwork.Nodes {
				fmt.Printf("\nHostname: %s", node.Hostname)
			}
			conn.Close()
		// Lists the shared files on a remote node.
		case "listFiles":
			switch req.State {
			case "new":
				node, err := mynetwork.FindNode(req.Target)
				if err != nil {
					fmt.Printf("`%s` is not in your network.", req.Target)
					conn.Close()
					return
				}
				listReq := network.Request{
					Target:  req.Target,
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Command: "listFiles",
					State:   "request",
				}
				node.SendRequest(listReq)
				conn.Close()
			case "request":
				node, err := mynetwork.FindNode(req.Source)
				if err != nil {
					fmt.Printf("`%s` is not in your network.", req.Source)
					conn.Close()
				}
				files := network.ListFiles(config.SharedDirectory)
				lfReq := network.Request{
					Target:  req.Source,
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Command: "listFiles",
					State:   "response",
					Body:    strings.Join(files, ","),
					Success: true,
				}
				node.SendRequest(lfReq)

			case "response":
				files := strings.Split(req.Body, ",")
				fmt.Printf("\n`%s` - Shared Files \n", req.Source)
				for _, file := range files {
					fmt.Println(file)
				}
			}
		case "download":
			switch req.State {
			case "new":
				node, err := mynetwork.FindNode(req.Target)
				if err != nil {
					fmt.Println("Unable to find Node!")
					conn.Close()
				}
				shared_dir := config.SharedDirectory
				dReq := network.Request{
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Target:  req.Target,
					Command: "download",
					State:   "response",
					Args:    req.Args,
				}
				fmt.Printf("\nSending using Conn: %s", node.Conn.LocalAddr())
				network.Download(node, dReq, shared_dir)
				conn.Close()
			case "response":
				// node, err := mynetwork.FindNode(req.Source)
				// if err != nil {
				// 	fmt.Println("Unable to find Node!")
				// 	return
				// }
				fileName := strings.TrimSpace(req.Args)
				shared_dir := config.SharedDirectory
				path_to_file := fmt.Sprintf("%s/%s", shared_dir, fileName)
				fmt.Printf("\nSending Response using Conn: %s", conn.RemoteAddr())
				network.SendFileToClient(conn, req, path_to_file)
				// conn.Close()
			}
		case "ping":
			for _, node := range mynetwork.Nodes {
				pingReq := network.Request{
					Target:  node.Hostname,
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Command: "pong",
					State:   "new",
				}
				node.SendRequest(pingReq)
			}
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
