package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"your_network/configuration"
)

const NETWORKFILE = "./network.json"

type Network struct {
	Nodes  []Node                      `json:"nodes"`
	Config configuration.Configuration `json:"-"`
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
		req.BlockingSend()
		// conn, err := net.Dial("tcp4", node.Hostname)
		// if err != nil {
		// 	fmt.Printf("\nFailed to establish connection with %s", node.Hostname)
		// 	os.Exit(1)
		// }
		// fmt.Printf("\nEstablish connection with %s - %s", node.Identifier, node.Hostname)
		// n.Nodes[i].Conn = conn
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
		fmt.Printf("`%s` == `%s`", nNode.Hostname, target)
		if nNode.Hostname == target {
			return true
		}
	}
	return false
}

func (n Network) AddNodeOnExisting(conn net.Conn, hostname string) *Network {
	uri := strings.TrimSpace(hostname)
	node := Node{
		Hostname: uri,
		Conn:     conn,
	}
	n.Nodes = append(n.Nodes, node)
	n.Update()
	return &n
}

func (n Network) AddNode(hostname string) *Network {
	uri := strings.TrimSpace(hostname)
	conn, err := net.Dial("tcp4", uri)
	if err != nil {
		fmt.Printf("\nFailed to establish connection with %s", uri)
		os.Exit(1)
	}
	node := Node{
		Hostname: hostname,
		Conn:     conn,
	}
	n.Nodes = append(n.Nodes, node)
	n.Update()
	return &n
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
