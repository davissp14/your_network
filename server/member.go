package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Member struct {
	Hostname string
	Conn     net.Conn
}

func NewMember(host, port string) (Member, error) {
	uri := strings.TrimSpace(fmt.Sprintf("%s:%s", host, port))
	conn, err := net.Dial("tcp4", uri)
	if err != nil {
		return Member{}, err
	}
	return Member{Hostname: uri, Conn: conn}, nil
}

func (m Member) Message(message string) {
	writer := bufio.NewWriter(m.Conn)
	fmt.Fprintln(writer, message)
	err := writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush write to %s: Error: %s", m.Hostname, err)
	}
}
