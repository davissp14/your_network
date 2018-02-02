package client

import (
	"fmt"
	"os"
	"your_network/configuration"
	"your_network/network"
	"your_network/server"
)

type Client struct {
	Config configuration.Configuration
}

func (c Client) StartServer() {
	server.Start(c.Config)
}

func (c Client) AddNode(target string) {
	if c.Self() == target {
		fmt.Println("Why are you trying to add yourself?")
		os.Exit(1)
	}
	req := network.Request{
		Target:        c.Self(),
		Source:        c.Self(),
		CommandTarget: target,
		Command:       "membership",
		State:         "request",
		Client:        true,
	}
	conn, err := network.NewConnection(c.Self())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Send(req)
	res := conn.Receive()
	fmt.Println(res.Body)
}

func (c Client) ListNodes() {
	req := network.Request{
		Target:  c.Self(),
		Source:  c.Self(),
		Command: "list_nodes",
		Client:  true,
	}
	conn, err := network.NewConnection(c.Self())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Send(req)
	res := conn.Receive()
	fmt.Println(res.Body)
}

//
// func (c Client) ListFiles(hostname string) {
// 	listReq := network.Request{
// 		Target:        c.Self(),
// 		CommandTarget: hostname,
// 		Command:       "listFiles",
// 		State:         "request",
// 		Client:        true,
// 	}
// 	conn, err := listReq.BlockingSend()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	req := listReq.BlockingRead(conn)
// 	fmt.Println(req.Body)
// }
//
// func (c Client) Download(target string, filename string) {
// 	downloadReq := network.Request{
// 		Target:        c.Self(),
// 		CommandTarget: target,
// 		Command:       "download",
// 		State:         "request",
// 		Args:          filename,
// 		Client:        true,
// 	}
// 	conn, err := downloadReq.BlockingSend()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	req := downloadReq.BlockingRead(conn)
// 	fmt.Println(req.Body)
// }
//
// func (c Client) Ping(target string) {
// 	pingReq := network.Request{
// 		Target:        c.Self(),
// 		CommandTarget: target,
// 		Command:       "ping",
// 		State:         "request",
// 		Client:        true,
// 	}
// 	conn, err := pingReq.BlockingSend()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	req := pingReq.BlockingRead(conn)
// 	fmt.Println(req.Body)
// }

func (c Client) Self() string {
	return fmt.Sprintf("%s:%s", c.Config.Hostname, c.Config.Port)
}
