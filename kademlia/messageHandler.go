package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

func (network *Network) HandleMessage (rawMsg []byte, recieverAddr *net.UDPAddr) ([]byte, error){
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
	case "STORE":
		fmt.Println("Received STORE from ", msg.Sender)
	case "FIND_NODE":
		fmt.Println("Received FIND_NODE from ", msg.Sender)
	case "FIND_VALUE":
		fmt.Println("Received FIND_VALUE from ", msg.Sender)
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