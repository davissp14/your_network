package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Command string `json:"command"`
	Body    string `json:"body"`
	Success bool   `json:"success"`
}

func (r Request) String() string {
	reqJSON, _ := json.Marshal(r)
	return string(reqJSON)
}

func ParseRequest(raw string) (Request, error) {
	var req Request
	json.Unmarshal([]byte(raw), &req)
	return req, nil
}

func (r Request) Send() {
	json, _ := json.Marshal(r)
	conn, err := net.Dial("tcp4", r.Target)
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
