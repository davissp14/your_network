package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Node struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key"`
	Conn      net.Conn
}

func NewNode(target string) (Node, error) {
	uri := strings.TrimSpace(target)
	conn, err := net.Dial("tcp4", uri)
	if err != nil {
		return Node{}, err
	}
	return Node{Hostname: uri, Conn: conn}, nil
}

func (n Node) SendRequest(req Request) {
	json, _ := json.Marshal(req)
	writer := bufio.NewWriter(n.Conn)
	fmt.Fprintln(writer, string(json))
	err := writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush write to %s: Error: %s", n.Hostname, err)
	}
}
