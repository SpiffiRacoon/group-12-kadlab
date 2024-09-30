package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNetwork(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	network1 := NewNetwork(contact1)

	contact2 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	network2 := NewNetwork(contact2)

	t.Run("Test SendPingMessage", func(t *testing.T) {
		go func (){
			err := network2.Listen("0.0.0.0", 8001)
			assert.Nil(t, err)
		}()

		//time.Sleep(1 * time.Second)

		response := network1.SendPingMessage(&contact2)
		assert.True(t, response)
	})

}

