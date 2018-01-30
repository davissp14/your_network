package server

import (
	"bufio"
	"encoding/json"
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
			case "request":
				response := network.Request{}
				if req.Source == req.CommandTarget {
					response.Success = false
					response.Body = "There's no need to add yourself."
					response.SendOnExisting(conn)
					conn.Close()
					return
				}
				memReq := network.Request{
					Target:  req.CommandTarget,
					Source:  req.Source,
					Command: "membershipRequest",
					State:   "response",
				}
				memConn, err := memReq.BlockingSend()
				if err != nil {
					response.Success = false
					response.Body = fmt.Sprintf("Failed to send request to %s", req.CommandTarget)
					response.SendOnExisting(conn)
					conn.Close()
					return
				}
				resReq := memReq.BlockingRead(memConn)
				if resReq.Success {
					mynetwork = mynetwork.AddNodeOnExisting(memConn, req.CommandTarget)
					fmt.Printf("`%s` joined your network!\n", req.CommandTarget)
					networkJSON, _ := json.Marshal(mynetwork)
					response.Success = true
					response.Body = string(networkJSON)
				} else {
					response.Success = false
					response.Body = resReq.Body
				}
				response.SendOnExisting(conn)
				conn.Close() // Close Client Connection
			case "response":
				response := network.Request{}
				if !mynetwork.NodeExists(req.Source) {
					mynetwork = mynetwork.AddNode(req.Source)
					fmt.Printf("`%s` joined your network!\n", req.Source)
					networkJSON, _ := json.Marshal(mynetwork)
					response.Success = true
					response.Body = string(networkJSON)
				} else {
					fmt.Printf("%s is already part of your network!\n", req.Source)
					response.Success = false
					response.Body = "Node is already part of your network!"
				}
				response.SendOnExisting(conn)
			}
			// Lists all of the nodes in the network.
		case "listNodes":
			networkJSON, _ := json.MarshalIndent(mynetwork, "", "  ")
			listNodes := network.Request{
				Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
				Success: true,
				Body:    string(networkJSON),
			}
			listNodes.SendOnExisting(conn)
			conn.Close() // Close Client connection

		case "ping":
			switch req.State {
			case "request":
				response := network.Request{}
				if !mynetwork.NodeExists(req.CommandTarget) {
					response.Success = false
					response.Body = fmt.Sprintf("%s isn't in your network.", req.CommandTarget)
					response.SendOnExisting(conn)
					conn.Close()
					return
				} else {
					pingReq := network.Request{
						Target:  req.CommandTarget,
						Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
						Command: "ping",
						State:   "response",
					}
					node, _ := mynetwork.FindNode(req.CommandTarget)
					pingReq.SendOnExisting(node.Conn)
					reqRes := pingReq.BlockingRead(node.Conn)
					response.Success = true
					response.Body = reqRes.Body
					response.SendOnExisting(conn)
					conn.Close()
				}
			case "response":
				resp := network.Request{
					Source:  req.Target,
					Body:    "PONG",
					Success: true,
				}
				resp.SendOnExisting(conn)
			}
		// Lists the shared files on a remote node.
		case "listFiles":
			switch req.State {
			case "request":
				response := network.Request{}
				if !mynetwork.NodeExists(req.CommandTarget) {
					response.Success = false
					response.Body = fmt.Sprintf("%s isn't in your network.", req.CommandTarget)
					response.SendOnExisting(conn)
					conn.Close()
				} else {
					node, _ := mynetwork.FindNode(req.CommandTarget)
					listReq := network.Request{
						Target:  req.CommandTarget,
						Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
						Command: "listFiles",
						State:   "response",
					}
					listReq.SendOnExisting(node.Conn)
					reqRes := listReq.BlockingRead(node.Conn)
					response.Success = true
					response.Body = reqRes.Body
					response.SendOnExisting(conn)
					conn.Close()
				}
			case "response":
				files := network.ListFiles(config.SharedDirectory)
				var shared_files network.SharedFileList
				for _, file := range files {
					shared_files.Files = append(shared_files.Files, file)
				}
				fileJSON, _ := json.MarshalIndent(shared_files, "", "  ")
				response := network.Request{
					Source:  req.Target,
					Success: true,
					Body:    string(fileJSON),
				}
				response.SendOnExisting(conn)
			}
		case "download":
			switch req.State {
			case "request":
				response := network.Request{}
				node, err := mynetwork.FindNode(req.CommandTarget)
				if err != nil {
					response.Body = "Node is not in your network"
					response.Success = false
					response.SendOnExisting(conn)
					conn.Close()
					return
				}
				dReq := network.Request{
					Source:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
					Target:  req.CommandTarget,
					Command: "download",
					State:   "response",
					Args:    req.Args,
				}
				fmt.Printf("\nSending using Conn: %s", node.Conn.LocalAddr())
				network.Download(node, dReq, config.SharedDirectory)
				response.Success = true
				response.Body = fmt.Sprintf("Successfully downloaded file `%s`", req.Args)
				response.SendOnExisting(conn)
				conn.Close()
			case "response":
				fileName := strings.TrimSpace(req.Args)
				shared_dir := config.SharedDirectory
				path_to_file := fmt.Sprintf("%s/%s", shared_dir, fileName)
				network.SendFileToClient(conn, req, path_to_file)
				// conn.Close()
			}
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
