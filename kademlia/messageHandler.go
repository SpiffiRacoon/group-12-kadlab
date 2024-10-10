package kademlia

import (
	"encoding/json"
	"errors"
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
		response := network.handlePingMessage()
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "JOIN":
		response := network.handleJoinMessage(msg.Sender)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "FIND_CONTACT":
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
		response := network.handleStoreMessage(msg.Content)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response")
			return nil, err
		}
		return responseBytes, nil
	case "FIND_VALUE":
		fmt.Println("Received FIND_VALUE from ", msg.Sender)
		target := msg.Content
		resp, err := network.handleFindDataMessage(target)
		return resp, err
	default:
		return nil, errors.New("Unknown message type: " + msg.MsgType)
	}
}

func (network *Network) handlePingMessage() Message {
	pong := Message{
		MsgType: "PONG",
		Content: "I'm alive",
		Sender:  network.Me,
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

func (network *Network) handleStoreMessage(content string) Message {
	splitContent := strings.Split(content, ";")
	_, exists := network.kademlia.ExtractData(splitContent[0])
	if exists {
		errMsg := Message{
			MsgType: "ERROR_STORE",
			Content: "Store location is occupied",
		}
		return errMsg
	} else {
		network.kademlia.LocalStorage([]byte(splitContent[1]), splitContent[0])
		storedMsg := Message{
			MsgType: "STORED",
			Content: "Data stored successfully on node",
		}
		return storedMsg
	}
}

func (network *Network) handleFindDataMessage(key string) ([]byte, error) {
	val, exists := network.kademlia.ExtractData(key)
	if exists {
		res, err := json.Marshal(string(val))
		return res, err

	} else {
		suggestedContacts := network.RoutingTable.FindClosestContacts(network.RoutingTable.me.ID, 5)
		res, err := json.Marshal(suggestedContacts)
		return res, err
	}
}
