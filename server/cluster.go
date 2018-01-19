package server

import (
	"bufio"
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

func (c Cluster) Communicate(message string) {
	for _, member := range c.Members {
		writer := bufio.NewWriter(member.Conn)
		fmt.Fprintln(writer, message)
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

func (c Cluster) MemberExists(member Member) bool {
	for _, m := range c.Members {
		if m.Hostname == member.Hostname {
			return true
		}
	}
	return false
}
