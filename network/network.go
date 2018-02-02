package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"your_network/configuration"
)

const NETWORKFILE = "./network.json"

type Network struct {
	Nodes  []Node                      `json:"nodes"`
	Config configuration.Configuration `json:"-"`
}

func (n *Network) HandleRequest(connection Connection, req Request) {
	switch req.Command {
	case "membership":
		req.Membership(connection, n)
	case "list_nodes":
		req.ListNodes(connection, n)
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

func (n Network) FindNode(target string) (Node, error) {
	for _, nNode := range n.Nodes {
		if nNode.Hostname == target {
			return nNode, nil
		}
	}
	return Node{}, errors.New("Could not find node.")
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

//
// func (c Cluster) MemberHosts() []string {
// 	var list []string
// 	for _, member := range c.Members {
// 		list = append(list, member.Hostname)
// 	}
// 	return list
// }
//
// func (c Cluster) FindMember(hostname string) (Member, error) {
// 	for _, m := range c.Members {
// 		if m.Hostname == hostname {
// 			return m, nil
// 		}
// 	}
// 	return Member{}, errors.New("Member not found")
// }
//
// func (c Cluster) MemberExists(hostname string) bool {
// 	for _, m := range c.Members {
// 		if m.Hostname == hostname {
// 			return true
// 		}
// 	}
// 	return false
// }
