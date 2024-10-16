package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	contact1 := NewContact(NewKademliaID("1000000000000000000000000000000000000000"), "localhost:8001")
	node1 := NewKademlia(contact1, false)
	network1 := NewNetwork(contact1, node1)

	contact2 := NewContact(NewKademliaID("2000000000000000000000000000000000000000"), "localhost:8002")
	node2 := NewKademlia(contact2, false)
	network2 := NewNetwork(contact2, node2)

	contact3 := NewContact(NewKademliaID("3000000000000000000000000000000000000000"), "localhost:8003")

	t.Run("Test Listen", func(t *testing.T) {
		go func() {
			err := network2.Listen("0.0.0.0", 8002)
			assert.Nil(t, err)
		}()
	})

	t.Run("Test SendPingMessage", func(t *testing.T) {
		err := network1.SendPingMessage(&contact2)
		assert.Nil(t, err)

		err = network1.SendPingMessage(&contact3)
		assert.NotNil(t, err)
		//assert.Equal(t, Error()
	})

	t.Run("Test SendJoinMessage", func(t *testing.T) {
		//Result: contact1 is added to contact2's routing table
		err := network1.SendJoinMessage(&contact2)
		assert.Nil(t, err)

		//Result: contact3 is not added to contact2's routing table
		err = network1.SendJoinMessage(&contact3)
		assert.NotNil(t, err)
	})

	t.Run("Test SendFindContactMessage", func(t *testing.T) {
		//Result: contact1 is found in contact2's routing table
		contacts, err := network1.SendFindContactMessage(&contact2, contact1.ID)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(contacts))
		assert.Equal(t, contact1, contacts[0])

		//Result: contact1 is not found in contact3's routing table
		contacts, err = network1.SendFindContactMessage(&contact3, contact1.ID)
		assert.NotNil(t, err)
		assert.Equal(t, 0, len(contacts))

		//Result: contact1 is found in contact2's routing table
		contacts, err = network1.SendFindContactMessage(&contact2, contact3.ID)
		assert.Nil(t, err)
		assert.NotContains(t, contacts, contact3)
	})

	t.Run("Test SendStoreMessage", func(t *testing.T) {
		//Result: data is stored in contact2's routing table
		key := node1.MakeKey([]byte("data"))
		response := network1.SendStoreMessage([]byte("data"), key, &contact2)
		assert.Nil(t, response)
	})

	t.Run("Test SendFindDataMessage", func(t *testing.T) {
		key := node1.MakeKey([]byte("data"))
		response, _ := network1.SendFindDataMessage(key, &contact2)
		assert.Equal(t, "data", response)
	})

}
