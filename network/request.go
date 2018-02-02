package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Source        string `json:"source"`
	Target        string `json:"target"`
	CommandTarget string `json:"command_target"`
	Command       string `json:"command"`
	Args          string `json:"args"`
	State         string `json:"state"`
	Body          string `json:"body"`
	Success       bool   `json:"success"`
	Client        bool   `json:"client"`
}

func (r Request) String() string {
	reqJSON, _ := json.Marshal(r)
	return string(reqJSON)
}

func ParseRequest(raw string) (Request, error) {
	var req Request
	json.Unmarshal([]byte(raw), &req)
	if req.State == "" {
		req.State = "new"
	}
	return req, nil
}

func (r Request) BlockingSend() (net.Conn, error) {
	json, _ := json.Marshal(r)
	conn, err := net.Dial("tcp", r.Target)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, string(json))
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (r Request) BlockingRead(conn net.Conn) Request {
	var req Request
	reader := bufio.NewReader(conn)
	for {
		data, _ := reader.ReadString('\n')
		if data == "" {
			break
		}
		json.Unmarshal([]byte(data), &req)
		return req
	}
	return req
}

func (r Request) SendOnExisting(conn net.Conn) error {
	json, _ := json.Marshal(r)
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, string(json))
	err := writer.Flush()
	return err
}
