package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"your_network/network"
)

type Command struct {
	Network *network.Network
	Request network.Request
	Conn    net.Conn
}

func NewCommand(n *network.Network, req network.Request, conn net.Conn) Command {
	return Command{
		Network: n,
		Request: req,
		Conn:    conn,
	}
}

func (c Command) UpdateNetwork() *network.Network {
	return c.Network
}

func (c Command) Handle() {
	switch c.Request.Command {
	case "membership":
		c.Membership()
	}
}

func (c Command) Membership() {
	switch c.Request.State {
	case "request":
		req := network.Request{Command: "membership", State: "response"}
		// Ensure node doesn't already exist in network.
		fmt.Println(len(c.Network.Nodes))
		if c.Network.NodeExists(c.Request.CommandTarget) {
			fmt.Println("HERE")
			c.SendResponse(c.Conn, "Node already exists in network", false)
			return
		}
		// Initiate connection with target node.
		nodeConn, err := net.Dial("tcp", c.Request.CommandTarget)
		if err != nil {
			c.SendResponse(c.Conn, "Failed to establish connection.", false)
			return
		}
		// Attempt to establish connection.
		c.Send(req, nodeConn)
		// Wait for response from target node.
		res := c.Receive(nodeConn)
		if !res.Success {
			c.SendResponse(c.Conn, res.Body, false)
			return
		}
		// Add node to network.
		c.Network = c.Network.AddNodeOnExisting(nodeConn, c.Request.CommandTarget)
		networkJSON, _ := json.Marshal(c.Network)
		c.SendResponse(c.Conn, string(networkJSON), true)

	case "response":
		if c.Network.NodeExists(c.Request.Source) {
			c.SendResponse(c.Conn, "Node already exists in network", false)
			return
		}
		c.Network = c.Network.AddNode(c.Request.Source)
		networkJSON, _ := json.Marshal(c.Network)
		c.SendResponse(c.Conn, string(networkJSON), true)
	}
}

func (c Command) Send(req network.Request, conn net.Conn) error {
	req.Target = c.Request.CommandTarget
	req.Source = c.Source()
	err := req.SendOnExisting(conn)
	if err != nil {
		return err
	}
	return nil
}

func (c Command) Receive(conn net.Conn) network.Request {
	var res network.Request
	reader := bufio.NewReader(conn)
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

func (c Command) SendResponse(conn net.Conn, body string, success bool) {
	resp := network.Request{
		Body:    body,
		Success: success,
	}
	resp.SendOnExisting(conn)
}

func (c Command) Source() string {
	return fmt.Sprintf("%s:%s", c.Network.Config.Hostname, c.Network.Config.Port)
}

// switch req.State {
// case "request":
//   response := network.Request{}
//   memReq := network.Request{
//     Target:  req.CommandTarget,
//     Source:  Source(mynetwork),
//     Command: "membership",
//     State:   "response",
//   }
//   memConn, err := memReq.BlockingSend()
//   if err != nil {
//     response.Success = false
//     response.Body = fmt.Sprintf("Failed to send request to %s", req.CommandTarget)
//     response.SendOnExisting(conn)
//     conn.Close()
//     return
//   }
//   resReq := memReq.BlockingRead(memConn)
//   if resReq.Success {
//     mynetwork = mynetwork.AddNodeOnExisting(memConn, req.CommandTarget)
//     fmt.Printf("`%s` joined your network!\n", req.CommandTarget)
//     // node, _ := mynetwork.FindNode(req.CommandTarget)
//     // go node.Monitor(Source(mynetwork))
//     networkJSON, _ := json.Marshal(mynetwork)
//     response.Success = true
//     response.Body = string(networkJSON)
//   } else {
//     response.Success = false
//     response.Body = resReq.Body
//   }
//   response.SendOnExisting(conn)
//   conn.Close() // Close Client Connection
// case "response":
//   response := network.Request{}
//   if !mynetwork.NodeExists(req.Source) {
//     mynetwork = mynetwork.AddNode(req.Source)
//     fmt.Printf("`%s` joined your network!\n", req.Source)
//     // node, err := mynetwork.FindNode(req.Source)
//     // go node.Monitor(Source(mynetwork))
//     networkJSON, _ := json.Marshal(mynetwork)
//     response.Success = true
//     response.Body = string(networkJSON)
//   } else {
//     fmt.Printf("%s is already part of your network!\n", req.Source)
//     response.Success = false
//     response.Body = "Node is already part of your network!"
//   }
//   response.SendOnExisting(conn)
// }
