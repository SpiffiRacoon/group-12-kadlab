package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	MsgType string
	Content string
	Sender  Contact
}

type Network struct {
	RoutingTable RoutingTable
	Me           Contact
}

func NewNetwork(me Contact) *Network {
	network := &Network{}
	network.Me = me
	network.RoutingTable = *NewRoutingTable(me)
	return network
}

func (network *Network) Listen(ip string, port int) error {
	lAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &lAddr)
	if err != nil {
		fmt.Println("Error listening on ", ip, ":", port)
		return err
	}
	defer conn.Close()
	fmt.Println("Listening on ", ip, ":", port)

	for {
		buffer := make([]byte, 1024)
		byteNum, rAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection")
			return err
		}
		response, err := network.HandleMessage(buffer[:byteNum], rAddr)
		if err != nil {
			fmt.Println("Error when handling Message:", err)
			continue
		}
		conn.WriteToUDP(response, rAddr)
	}
}

func (network *Network) SendMessage(msg Message, contact *Contact) ([]byte, error) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "not responding", err.Error())
		return nil, err
	}
	data, _ := json.Marshal(msg)
	_, err = conn.Write(data)

	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "sent error while writing", err.Error())
		return nil, err
	}

	deadline := time.Now().Add(15 * time.Second)
	conn.SetDeadline(deadline)

	response := make([]byte, 1024)
	byteNum, err := conn.Read(response)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "sent error while reading", err.Error())
		return nil, err
	}

	defer conn.Close()
	return response[:byteNum], err
}

func (network *Network) SendPingMessage(contact *Contact) bool {
	msg := Message{
		MsgType: "PING",
		Content: "PING",
		Sender:  *contact,
	}
	responseMsg, err := network.SendMessage(msg, contact)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "not responding", err.Error())
		return false
	} else {
		var msg Message
		err := json.Unmarshal(responseMsg, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message")
			return false
		}
		fmt.Printf("%s %s %s %s %s\n", contact.ID, "responding on port:", contact.Address, "with ", msg.Content)
		return true
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error) {
	// TODO
	return nil, nil
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, key string, contact *Contact) bool { //förslag, använd error istället för bools

	msg := Message{
		MsgType: "STORE",
		Content: key + ";" + string(data), //Order here can be reversed if needed but should not matter as long as you know the order
		Sender:  *contact,
	}
	responseMsg, err := network.SendMessage(msg, contact)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "not responding", err.Error())
		return false
	} else {
		var storeResponse Message
		err := json.Unmarshal(responseMsg, &storeResponse)
		if err != nil {
			fmt.Println("Error during unmarshalling")
			return false
		}
		if storeResponse.MsgType != "STORED" {
			fmt.Println("Failed to store")
			return false
		}
		return true
	}

}
