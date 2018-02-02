package network

import (
	"encoding/json"
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
		networkJSON, _ := json.Marshal(network)
		incoming.SendResponse(string(networkJSON), true)
		return
	default:
		incoming.SendResponse("You shouldn't be here!", false)
	}
}

func (r Request) ListNodes(incoming Connection, network *Network) {
	networkJSON, _ := json.MarshalIndent(network, "", "  ")
	incoming.SendResponse(string(networkJSON), true)
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
