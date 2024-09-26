package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

func (network *Network) HandleMessage(rawMsg []byte, recieverAddr *net.UDPAddr) ([]byte, error) {
	var msg Message
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling message")
		return nil, err
	}

	switch msg.MsgType {
	case "PING":
		fmt.Println("Received PING from ", msg.Sender)
		response := network.handlePingMessage()
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "JOIN":
		fmt.Println("Received JOIN from ", msg.Sender)
		response := network.handleJoinMessage(msg.Sender)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "FIND_CONTACT":
		fmt.Println("Received FIND_CONTACT from ", msg.Sender)
		target := NewKademliaID(msg.Content)
		contacts := network.handleFindContactMessage(target, 3)
		contactsBytes, err := json.Marshal(contacts)
		if err != nil {
			fmt.Println("Error marshalling contacts")
			return nil, err
		}
		return contactsBytes, nil
	case "STORE":
		fmt.Println("Received STORE from ", msg.Sender)
	case "FIND_VALUE":
		fmt.Println("Received FIND_VALUE from ", msg.Sender)
	default:
		fmt.Println("Unknown message type: " + msg.MsgType)
	}
	return nil, nil
}

func (network *Network) handlePingMessage() Message {
	pong := Message{
		MsgType: "PONG",
		Content: "I'm alive",
	}
	return pong
}

func (network *Network) handleJoinMessage(sender Contact) Message {
	network.RoutingTable.AddContact(sender)
	joinResponse := Message{
		MsgType: "JOIN_RESPONSE",
		Content: "Welcome to the network",
	}
	return joinResponse
}

func (network *Network) handleFindContactMessage(target *KademliaID, count int) []Contact {
	contacts := network.RoutingTable.FindClosestContacts(target, count)
	return contacts
}

func (network *Network) handleStoreMessage() Message {
	storedmsg := Message{
		MsgType: "STORED",
		Content: "Data stored successfully on node",
	}
	return storedmsg
}
