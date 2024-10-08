package kademlia

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKademlia(t *testing.T) {
	bootstrapContact := NewContact(NewKademliaID("B0075712A9000000000000000000000000000000"), "localhost:3100")
	contact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:3101")

	bootstrap := NewKademlia(bootstrapContact, true)
	assert.Equal(t, bootstrap.Me.ID, bootstrapContact.ID)
	assert.Equal(t, bootstrap.Me.Address, bootstrapContact.Address)
	assert.Nil(t, bootstrap.BootstrapNode.ID)
	assert.True(t, bootstrap.IsBootstrap)
	go bootstrap.Start()

	node := NewKademlia(contact, false)
	assert.Equal(t, node.Me.ID, contact.ID)
	assert.Equal(t, node.Me.Address, contact.Address)
	assert.Nil(t, node.BootstrapNode.ID)
	assert.False(t, node.IsBootstrap)

	node.BootstrapNode = bootstrapContact
	assert.Equal(t, node.BootstrapNode, bootstrapContact)

	go node.Start()

	t.Run("Test LookupContact", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		contacts, err := node.LookupContact(NewKademliaID("B0075712A9000000000000000000000000000000"))
		assert.Nil(t, err)
		assert.Equal(t, 3, len(contacts))
		fmt.Println(contacts)
	})

	t.Run("Test Store", func(t *testing.T) {
		stored := node.Store([]byte("TestingTesting"))
		assert.Nil(t, stored)
	})

	t.Run("Test LookupData", func(t *testing.T) {

		dataRes := node.LookupData(node.MakeKey([]byte("TestingTesting")))
		fmt.Print(string(dataRes))
		assert.NotNil(t, dataRes)
		anotherDataRes := node.LookupData(node.MakeKey([]byte("SkaInteFinnas")))
		assert.Nil(t, anotherDataRes)
	})
}
