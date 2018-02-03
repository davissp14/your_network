package network

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (r Request) Membership(incoming Connection, network *Network) {
	switch r.State {
	case "request":
		req := Request{
			Source:  r.Source,
			Target:  r.CommandTarget,
			Command: "membership",
			State:   "response",
		}
		// Ensure node doesn't already exist in network.
		if network.NodeExists(r.CommandTarget) {
			incoming.SendResponse("Node already exists in network", false)
			return
		}
		// Initiate connection with target node.
		outgoing, err := NewConnection(r.CommandTarget)
		if err != nil {
			incoming.SendResponse("Failed to establish connection.", false)
			return
		}
		// Attempt to establish connection.
		outgoing.Send(req)
		// Wait for response from target node.
		res := outgoing.Receive()
		if !res.Success {
			incoming.SendResponse(res.Body, false)
			return
		}
		// Add node to network.
		network.AddNode(outgoing)
		networkJSON, _ := json.Marshal(network)
		incoming.SendResponse(string(networkJSON), true)
	case "response":
		if network.NodeExists(r.Source) {
			incoming.SendResponse("Node already exists in network", false)
			return
		}
		outgoing, _ := NewConnection(r.Source)
		network.AddNode(outgoing)
		networkJSON, _ := json.MarshalIndent(network, "", "  ")
		incoming.SendResponse(string(networkJSON), true)
		return
	}
}

func (r Request) ListNodes(incoming Connection, network *Network) {
	networkJSON, _ := json.MarshalIndent(network, "", "  ")
	incoming.SendResponse(string(networkJSON), true)
}

func (r Request) ListFiles(incoming Connection, network *Network) {
	switch r.State {
	case "request":
		if !network.NodeExists(r.CommandTarget) {
			incoming.SendResponse("Node isn't in your network", false)
			return
		}
		req := Request{
			Target:  r.CommandTarget,
			Source:  r.Source,
			Command: "list_files",
			State:   "response",
		}
		outgoing, err := network.FindConnection(r.CommandTarget)
		if err != nil {
			incoming.SendResponse("Failed to find the connection", false)
			return
		}
		outgoing.Send(req)
		res := outgoing.Receive()
		if !res.Success {
			incoming.SendResponse(res.Body, false)
			return
		}
		incoming.SendResponse(res.Body, true)
	case "response":
		files := ListFiles(network.Config.SharedDirectory)
		var sharedFiles SharedFileList
		for _, file := range files {
			sharedFiles.Files = append(sharedFiles.Files, file)
		}
		fileJSON, _ := json.MarshalIndent(sharedFiles, "", "  ")
		req := Request{
			Source:  r.Target,
			Success: true,
			Body:    string(fileJSON),
		}
		incoming.Send(req)
	}
}

func (r Request) Download(incoming Connection, network *Network) {
	switch r.State {
	case "request":
		if !network.NodeExists(r.CommandTarget) {
			incoming.SendResponse("Node isn't in your network", false)
			return
		}
		req := Request{
			Source:  r.Source,
			Target:  r.CommandTarget,
			Command: "download",
			State:   "response",
			Args:    r.Args,
		}
		outgoing, _ := network.FindConnection(r.CommandTarget)
		fmt.Printf("\nSending using Conn: %s", outgoing.Conn.LocalAddr())
		Download(outgoing, req, network.Config.SharedDirectory)
		incoming.SendResponse("Successfully downloaded file!", true)
	case "response":
		fileName := strings.TrimSpace(r.Args)
		sharedDir := network.Config.SharedDirectory
		pathToFile := fmt.Sprintf("%s/%s", sharedDir, fileName)
		SendFileToClient(incoming, pathToFile)
	}
}

func (r Request) Ping(incoming Connection, network *Network) {
	switch r.State {
	case "request":
		if !network.NodeExists(r.CommandTarget) {
			incoming.SendResponse("Node isn't in your network", false)
			return
		}
		req := Request{
			Target:  r.CommandTarget,
			Source:  r.Source,
			Command: "ping",
			State:   "response",
		}
		outgoing, _ := network.FindConnection(r.CommandTarget)
		outgoing.Send(req)
		res := outgoing.Receive()
		if !res.Success {
			incoming.SendResponse(res.Body, false)
			return
		}
		incoming.SendResponse(res.Body, true)
	case "response":
		req := Request{
			Source:  r.Target,
			Body:    "PONG",
			Success: true,
		}
		incoming.Send(req)
	}
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
