package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
	RoutingTable RoutingTable
	Me 	   Contact
}

type Message struct {
	MsgType string
	Content   string
	Sender Contact
}

func NewNetwork(me Contact) *Network {
	network := &Network{}
	network.Me = me
	network.RoutingTable = *NewRoutingTable(me)
	return network
}

func (network *Network) Listen(ip string, port int) error{
	// TODO
	return nil
}

func (network *Network) SendPingMessage(contact *Contact) bool{
	//timeout := time.Duration(1 * time.Second)
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "not responding", err.Error())
		return false
	} else {
		fmt.Printf("%s %s %s\n", contact.ID, "responding on port:", contact.Address)
	}
	defer conn.Close()
	return true
}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error) {
	// TODO
	return nil, nil
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
