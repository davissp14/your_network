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
				memReq := network.Request{
					Target:  req.CommandTarget,
					Source:  Source(mynetwork),
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
					// node, _ := mynetwork.FindNode(req.CommandTarget)
					// go node.Monitor(Source(mynetwork))
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
					// node, err := mynetwork.FindNode(req.Source)
					// go node.Monitor(Source(mynetwork))
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
				Source:  Source(mynetwork),
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
						Source:  Source(mynetwork),
						Command: "ping",
						State:   "response",
					}
					node, _ := mynetwork.FindNode(req.CommandTarget)
					pingReq.SendOnExisting(node.Conn)
					reqRes := pingReq.BlockingRead(node.Conn)
					response.Success = true
					response.Body = reqRes.Body
					response.SendOnExisting(conn)
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
						Source:  Source(mynetwork),
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
				files := network.ListFiles(mynetwork.Config.SharedDirectory)
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
					Source:  Source(mynetwork),
					Target:  req.CommandTarget,
					Command: "download",
					State:   "response",
					Args:    req.Args,
				}
				fmt.Printf("\nSending using Conn: %s", node.Conn.LocalAddr())
				network.Download(node, dReq, mynetwork.Config.SharedDirectory)
				response.Success = true
				response.Body = fmt.Sprintf("Successfully downloaded file `%s`", req.Args)
				response.SendOnExisting(conn)
				conn.Close()
			case "response":
				fileName := strings.TrimSpace(req.Args)
				shared_dir := mynetwork.Config.SharedDirectory
				path_to_file := fmt.Sprintf("%s/%s", shared_dir, fileName)
				network.SendFileToClient(conn, req, path_to_file)
			}
		default:
			fmt.Printf("DAta: `%s`", req.String())
			conn.Close()
		}
	}
}

func Source(n network.Network) string {
	return fmt.Sprintf("%s:%s", n.Config.Hostname, n.Config.Port)
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
