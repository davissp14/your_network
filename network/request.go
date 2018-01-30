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

func (r Request) Send() {
	json, _ := json.Marshal(r)
	conn, err := net.Dial("tcp", r.Target)
	if err != nil {
		fmt.Println(err)
	}
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, string(json))
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush write to %s: Error: %s", r.Target, err)
	}
	conn.Close()
}

func (r Request) SendOnExisting(conn net.Conn) {
	json, _ := json.Marshal(r)
	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, string(json))
	err := writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush write to %s: Error: %s", r.Target, err)
	}
}
