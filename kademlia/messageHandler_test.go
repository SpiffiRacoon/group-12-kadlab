package kademlia

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	// Create a new network
	sender := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	target := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8002")
	target_node := NewKademlia(target, false)
	network := NewNetwork(target, target_node)

	// Create a new message

	t.Run("Test PING message", func(t *testing.T) {
		msg := Message{
			MsgType: "PING",
			Sender:  sender,
		}
		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.Nil(t, err)
		assert.NotNil(t, respData)

		var respMsg Message
		err = json.Unmarshal(respData, &respMsg)
		assert.Nil(t, err)
		assert.Equal(t, "PONG", respMsg.MsgType)
	})

	t.Run("Test JOIN message", func(t *testing.T) {
		msg := Message{
			MsgType: "JOIN",
			Content: sender.ID.String(),
			Sender:  sender,
		}

		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.Nil(t, err)
		assert.NotNil(t, respData)

		var respMsg Message
		err = json.Unmarshal(respData, &respMsg)
		assert.Nil(t, err)
		assert.Equal(t, "JOIN_RESPONSE", respMsg.MsgType)
	})

	t.Run("Test FIND_CONTACT message", func(t *testing.T) {
		network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002")) //000
		network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000001"), "localhost:8002")) //001
		network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000002"), "localhost:8002")) //010
		network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000003"), "localhost:8002")) //011
		network.RoutingTable.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000004"), "localhost:8002")) //100

		msg := Message{
			MsgType: "FIND_CONTACT",
			Content: "1111111100000000000000000000000000000002",
			Sender:  sender,
		}
		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.Nil(t, err)
		assert.NotNil(t, respData)
		var respContacts []Contact
		err = json.Unmarshal(respData, &respContacts)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(respContacts))
		assert.Contains(t, respContacts, NewContact(NewKademliaID("1111111100000000000000000000000000000002"), "localhost:8002"))
	})

	t.Run("Test STORE message", func(t *testing.T) {
		msg := Message{
			MsgType: "STORE",
			Content: network.kademlia.MakeKey([]byte("Hello World")) + ";" + "Hello World",
			Sender:  sender,
		}
		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.Nil(t, err)
		assert.NotNil(t, respData)
	})

	t.Run("Test FIND_VALUE message", func(t *testing.T) {
		msg := Message{
			MsgType: "FIND_VALUE",
			Content: network.kademlia.MakeKey([]byte("Hello World")),
			Sender:  sender,
		}
		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.Nil(t, err)
		assert.Equal(t, []byte("\"Hello World\""), respData)
		//TODO, fixa så messageHandlern returnerar utan backslashes(läggs till i samband med marshallingen)
		//Fixa assert för att kolla antalet kontakter som skickas tillbaka om värdet ej finns i nuvarande nod
		//Fundera på hur vi vill att kontakter skickas tillbaka, är det vettigt med bytearray eller måste det vara en contact array?
	})

	t.Run("Test unknown message", func(t *testing.T) {
		msg := Message{
			MsgType: "UNKNOWN",
			Sender:  sender,
		}
		data, _ := json.Marshal(msg)
		respData, err := network.HandleMessage(data, nil)
		assert.NotNil(t, err)
		assert.Equal(t, "Unknown message type: UNKNOWN", err.Error())
		assert.Nil(t, respData)
	})

}
