package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
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
		response := network.HandlePingMessage()
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "JOIN":
		fmt.Println("Received JOIN from ", msg.Sender)
		response := network.HandleJoinMessage(msg.Sender)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "FIND_CONTACT":
		fmt.Println("Received FIND_CONTACT from ", msg.Sender)
		target := NewKademliaID(msg.Content)
		contacts := network.HandleFindContactMessage(target, 3)
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
		target := msg.Content
		data, contacts := network.HandleFindDataMessage(target)
		if data == nil {
			contactsBytes, err := json.Marshal(contacts)
			if err != nil {
				fmt.Println("Error marshalling contacts")
				return nil, err
			}
			return contactsBytes, nil
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling data")
			return nil, err
		}
		return dataBytes, nil
	default:
		fmt.Println("Unknown message type: " + msg.MsgType)
	}
	return nil, nil
}

func (network *Network) HandlePingMessage() Message {
	pong := Message{
		MsgType: "PONG",
		Content: "I'm alive",
	}
	return pong
}

func (network *Network) HandleJoinMessage(sender Contact) Message {
	network.RoutingTable.AddContact(sender)
	joinResponse := Message{
		MsgType: "JOIN_RESPONSE",
		Content: "Welcome to the network",
	}
	return joinResponse
}

func (network *Network) HandleFindContactMessage(target *KademliaID, count int) []Contact {
	contacts := network.RoutingTable.FindClosestContacts(target, count)
	return contacts
}

func (network *Network) HandleStoreMessage(content string) Message {
	splitContent := strings.Split(content, ";")
	tryStore := network.kademlia.ExtractData(splitContent[0])
	if tryStore != nil {
		errMsg := Message{
			MsgType: "ERROR_STORE",
			Content: "Store location is occupied",
		}
		return errMsg
	} else {
		storedMsg := Message{
			MsgType: "STORED",
			Content: "Data stored successfully on node",
		}
		return storedMsg
	}
}

func (network *Network) HandleFindDataMessage(content string) ([]byte, []Contact) {
	splitContent := strings.Split(content, ";")
	tryFind := network.kademlia.ExtractData(splitContent[0])
	if tryFind != nil {
		return tryFind, nil
	} else {
		suggestedContacts := network.RoutingTable.FindClosestContacts(network.RoutingTable.me.ID, 5)
		return nil, suggestedContacts
	}
}

//buskig trattel
