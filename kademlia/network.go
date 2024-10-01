package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Network struct {
	RoutingTable RoutingTable
	Me           Contact
}
type Message struct {
	MsgType string
	Content string
	Sender  Contact
	Target  Contact
}

func NewNetwork(me Contact) *Network {
	network := &Network{}
	network.Me = me
	network.RoutingTable = *NewRoutingTable(me)
	return network
}

func (network *Network) Listen(ip string, port int) error {
	localAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &localAddr)
	if err != nil {
		fmt.Println("Error listening on ", ip, ":", port)
		return err
	}
	defer conn.Close()
	fmt.Println("Listening on ", ip, ":", port)

	for {
		buffer := make([]byte, 1024)
		byteNum, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection")
			return err
		}
		response, err := network.HandleMessage(buffer[:byteNum], remoteAddr)
		if err != nil {
			fmt.Println("Error when handling Message:", err)
			continue
		}
		conn.WriteToUDP(response, remoteAddr)
	}
}

func (network *Network) sendMessage(msg Message, contact *Contact) ([]byte, error) {
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

func (network *Network) SendPingMessage(target *Contact) bool {
	msg := Message{
		MsgType: "PING",
		Content: "PING",
		Sender:  network.Me,
	}

	responseMsg, err := network.sendMessage(msg, target)
	if err != nil {
		fmt.Printf("%s %s %s\n", target.ID, "not responding", err.Error())
		return false
	} else {
		var msg Message
		err := json.Unmarshal(responseMsg, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message")
			return false
		}
		fmt.Printf("%s %s %s %s %s\n", target.ID, "responding on port:", target.Address, "with ", msg.Content)
		return true
	}
}

func (network *Network) SendJoinMessage(contact *Contact) bool {
	msg := Message{
		MsgType: "JOIN",
		Content: network.Me.ID.String(),
		Sender:  network.Me,
	}

	responseMsg, err := network.sendMessage(msg, contact)
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

func (network *Network) SendFindContactMessage(contact *Contact, targetID *KademliaID) ([]Contact, error) {
	msg := Message{
		MsgType: "FIND_CONTACT",
		Content: targetID.String(),
		Sender:  network.Me,
	}

	contactsByte, err := network.sendMessage(msg, contact)
	if err != nil {
		fmt.Printf("%s %s %s\n", contact.ID, "not responding", err.Error())
		return nil, err
	}

	var contacts []Contact
	err = json.Unmarshal(contactsByte, &contacts)
	if err != nil {
		fmt.Println("Error unmarshalling contacts")
		return nil, err
	}

	return contacts, nil

}

func (network *Network) SendStoreMessage(data []byte, key string, contact *Contact) bool { //förslag, använd error istället för bools

	msg := Message{
		MsgType: "STORE",
		Content: key + ";" + string(data), //Order here can be reversed if needed but should not matter as long as you know the order
		Sender:  network.RoutingTable.me,
		Target:  *contact,
	}
	responseMsg, err := network.sendMessage(msg, contact)
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

func (network *Network) SendFindDataMessage(hash string, contact *Contact) (string, bool) {
	msg := Message{
		MsgType: "FIND",
		Content: hash,
		Sender:  network.RoutingTable.me,
		Target:  *contact,
	}
	response, err := network.sendMessage(msg, contact)
	if err != nil {
		return ("Failed to contact target node."), false
	}
	var dataResponse string
	err = json.Unmarshal(response, &dataResponse)
	if err != nil {
		return "Error during unmarshalling", false
	} else if dataResponse == "" {
		return "Data not found", false
	} else {
		return dataResponse, true
	}
}


