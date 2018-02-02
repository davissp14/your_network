package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Connection struct {
	Conn     net.Conn
	Hostname string
}

func NewConnection(hostname string) (Connection, error) {
	hostname = strings.TrimSpace(hostname)
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		return Connection{}, err
	}
	return Connection{Conn: conn, Hostname: hostname}, nil
}

func (c Connection) Send(req Request) error {
	json, _ := json.Marshal(req)
	writer := bufio.NewWriter(c.Conn)
	fmt.Fprintln(writer, string(json))
	err := writer.Flush()
	return err
}

func (c Connection) Receive() Request {
	var res Request
	reader := bufio.NewReader(c.Conn)
	for {
		data, _ := reader.ReadString('\n')
		if data == "" {
			break
		}
		json.Unmarshal([]byte(data), &res)
		return res
	}
	return res
}

func (c Connection) SendResponse(body string, success bool) {
	resp := Request{
		Body:    body,
		Success: success,
	}
	c.Send(resp)
}
