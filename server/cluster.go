package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
)

type Message struct {
	Source  string
	Target  string
	Command string
	Success bool
}

type Cluster struct {
	Members []Member
}

func (c Cluster) addMember(member Member) Cluster {
	c.Members = append(c.Members, member)
	return c
}

func (c Cluster) Ping() {
	for _, member := range c.Members {
		req := Request{
			Target:  member.Hostname,
			Command: "pong",
		}
		json, _ := json.Marshal(req)
		writer := bufio.NewWriter(member.Conn)
		fmt.Fprintln(writer, string(json))
		err := writer.Flush()
		if err != nil {
			fmt.Printf("Failed to flush write to %s: Error: %s", member.Hostname, err)
		}
	}
}

func (c Cluster) MemberHosts() []string {
	var list []string
	for _, member := range c.Members {
		list = append(list, member.Hostname)
	}
	return list
}

func (c Cluster) FindMember(hostname string) (Member, error) {
	for _, m := range c.Members {
		if m.Hostname == hostname {
			return m, nil
		}
	}
	return Member{}, errors.New("Member not found")
}

func (c Cluster) MemberExists(hostname string) bool {
	for _, m := range c.Members {
		if m.Hostname == hostname {
			return true
		}
	}
	return false
}
