package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Member struct {
	Hostname string
	Conn     net.Conn
}

func NewMember(target string) (Member, error) {
	uri := strings.TrimSpace(target)
	conn, err := net.Dial("tcp4", uri)
	if err != nil {
		return Member{}, err
	}
	return Member{Hostname: uri, Conn: conn}, nil
}

func (m Member) SendRequest(req Request) {
	json, _ := json.Marshal(req)
	fmt.Println("Sending: ", string(json))
	writer := bufio.NewWriter(m.Conn)
	fmt.Fprintln(writer, string(json))
	err := writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush write to %s: Error: %s", m.Hostname, err)
	}
}
