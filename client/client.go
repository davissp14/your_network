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
		CommandTarget: target,
		Command:       "membership",
		State:         "request",
	}
	c.Send(req)
}

func (c Client) ListNodes() {
	req := network.Request{
		Command: "list_nodes",
	}
	c.Send(req)
}

func (c Client) ListFiles(hostname string) {
	req := network.Request{
		CommandTarget: hostname,
		Command:       "list_files",
		State:         "request",
	}
	c.Send(req)
}

func (c Client) Download(target string, filename string) {
	req := network.Request{
		CommandTarget: target,
		Command:       "download",
		State:         "request",
		Args:          filename,
	}
	c.Send(req)
}

func (c Client) Ping(target string) {
	req := network.Request{
		CommandTarget: target,
		Command:       "ping",
		State:         "request",
	}
	c.Send(req)
}

func (c Client) Send(req network.Request) {
	req.Source = c.Self()
	req.Client = true
	conn, err := network.NewConnection(c.Self())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Send(req)
	res := conn.Receive()
	fmt.Println(res.Body)
}

func (c Client) Self() string {
	return fmt.Sprintf("%s:%s", c.Config.Hostname, c.Config.Port)
}
