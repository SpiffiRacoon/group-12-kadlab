package kademlia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FIXME: This test doesn't actually test anything. There is only one assertion
// that is included as an example.

func TestRoutingTable(t *testing.T) {
	contact1 := NewContact(NewKademliaID("1000000000000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("2000000000000000000000000000000000000000"), "localhost:8000")
	contact3 := NewContact(NewKademliaID("3000000000000000000000000000000000000000"), "localhost:8000")
	contact4 := NewContact(NewKademliaID("4000000000000000000000000000000000000000"), "localhost:8000")
	contact5 := NewContact(NewKademliaID("5000000000000000000000000000000000000000"), "localhost:8000")
	contact6 := NewContact(NewKademliaID("6000000000000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(contact1)
	rt.AddContact(contact2)
	rt.AddContact(contact3)
	rt.AddContact(contact4)
	rt.AddContact(contact5)
	rt.AddContact(contact6)

	t.Run("Test FindClosestContacts", func(t *testing.T) {
		contacts := rt.FindClosestContacts(contact2.ID, 1)
		assert.Equal(t, 1, len(contacts))
		assert.Equal(t, contact2.ID, contacts[0].ID)

		contacts2 := rt.FindClosestContacts(contact2.ID, 3)
		assert.Equal(t, 3, len(contacts2))
		assert.Equal(t, contact2.ID, contacts2[0].ID)

		contacts3 := rt.FindClosestContacts(contact2.ID, 10)
		assert.Equal(t, 5, len(contacts3))
		assert.Equal(t, contact2.ID, contacts3[0].ID)

		for i := 1; i < len(contacts3); i++ {
			assert.True(t, contacts3[i-1].distance.Less(contacts3[i].distance))
		}

		for i := 0; i < len(contacts2); i++ {
			assert.Equal(t, contacts2[i], contacts3[i])
		}
	})

	t.Run("Test GenerateIDForBucket", func(t *testing.T) {
		for i := 0; i < IDLength*8; i=2*i+1 {
			id := rt.GenerateIDForBucket(i)
			bucket := rt.getBucketIndex(id)
			assert.Equal(t, i, bucket)
			fmt.Println("Generated ID: ", id)
		}
	})
}

