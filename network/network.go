package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"your_network/configuration"
)

type Network struct {
	Nodes  []Node                      `json:"nodes"`
	Config configuration.Configuration `json:"-"`
}

type Node struct {
	Hostname   string     `json:"hostname"`
	Alias      string     `json:"alias"`
	MacAddr    string     `json:"-"`
	PublicKey  string     `json:"-"`
	Connection Connection `json:"-"`
}

const NETWORKFILE = "./network.json"

func (n *Network) HandleRequest(connection Connection, req Request) {
	switch req.Command {
	case "membership":
		req.Membership(connection, n)
	case "list_nodes":
		req.ListNodes(connection, n)
	case "list_files":
		req.ListFiles(connection, n)
	case "download":
		req.Download(connection, n)
	case "ping":
		req.Ping(connection, n)
	}
}

func (n Network) Init() {
	if _, err := os.Stat(NETWORKFILE); os.IsNotExist(err) {
		fmt.Println("No network found. Initiating network.", NETWORKFILE)
		networkJSON, _ := json.Marshal(n)
		err := ioutil.WriteFile(NETWORKFILE, networkJSON, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Network established...")
	} else {
		fmt.Println("Loading existing network...")
		Load(n, n.Config)
	}
}

func Load(n Network, config configuration.Configuration) {
	raw, err := ioutil.ReadFile(NETWORKFILE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(raw, &n)
	for _, node := range n.Nodes {
		req := Request{
			Source:  node.Hostname,
			Target:  fmt.Sprintf("%s:%s", config.Hostname, config.Port),
			Command: "init_connection",
		}
		conn, _ := NewConnection(req.Target)
		conn.Send(req)
	}
}

func (n Network) Update() Network {
	networkJSON, _ := json.Marshal(n)
	err := ioutil.WriteFile(NETWORKFILE, networkJSON, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return n
}

func (n Network) FindConnection(target string) (Connection, error) {
	for _, node := range n.Nodes {
		if node.Hostname == target {
			return node.Connection, nil
		}
	}
	return Connection{}, errors.New("Could not find connection.")
}

func (n Network) NodeExists(target string) bool {
	for _, nNode := range n.Nodes {
		// fmt.Printf("`%s` == `%s`", nNode.Hostname, target)
		if nNode.Hostname == target {
			return true
		}
	}
	return false
}

func (n *Network) AddNode(connection Connection) {
	node := Node{
		Hostname:   connection.Hostname,
		Connection: connection,
	}
	n.Nodes = append(n.Nodes, node)
	n.Update()
}
