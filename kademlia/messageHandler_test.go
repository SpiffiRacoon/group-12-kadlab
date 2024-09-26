package kademlia

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	// Create a new network
	network := NewNetwork(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))	
	// Create a new message

	t.Run("Test PING message", func(t *testing.T) {
		msg := Message{
			MsgType: "PING",
			Sender: NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001"),
			
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

}