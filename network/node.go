package network

type Node struct {
	Hostname   string     `json:"hostname"`
	MacAddr    string     `json:"-"`
	PublicKey  string     `json:"-"`
	Connection Connection `json:"-"`
}

// func (n Node) Monitor(source string) {
// 	for {
// 		go func() {
// 			pingReq := Request{
// 				Source:        source,
// 				CommandTarget: source,
// 				Target:        n.Hostname,
// 				Command:       "ping",
// 				State:         "request",
// 			}
// 			pingReq.SendOnExisting(n.Conn)
// 			req := pingReq.BlockingRead(n.Conn)
// 			if req.Success {
// 				fmt.Printf("PONG - %s (latency)\n", n.Hostname)
// 			} else {
// 				fmt.Println(req.Body)
// 			}
// 		}()
// 		time.Sleep(5 * time.Second)
// 	}
// }
