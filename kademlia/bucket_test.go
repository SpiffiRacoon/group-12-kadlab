package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	// Create a new bucket
	contact0 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000")
	network := NewNetwork(contact0, nil)
	bucket := newBucket(network)

	// Create a new contact
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:8000")
	contact3 := NewContact(NewKademliaID("22222222000000000000000000000000000000000"), "localhost:8000")

	t.Run("Test AddContact", func(t *testing.T) {
		// Add a contact to the bucket
		bucket.AddContact(contact1)
		assert.Equal(t, 1, bucket.Len())

		bucket.AddContact(contact2)
		bucket.AddContact(contact3)
		assert.Equal(t, 3, bucket.Len())
	})

	t.Run("Test GetContactAndCalcDistance", func(t *testing.T) {
		// Get a contact from the bucket
		contacts := bucket.GetContactAndCalcDistance(NewKademliaID("0000000000000000000000000000000000000000"))

		assert.Equal(t, 3, len(contacts))

		assert.Equal(t, contacts[2].distance, contact1.ID)
		assert.Equal(t, contacts[1].distance, contact2.ID)
		assert.Equal(t, contacts[0].distance, contact3.ID)
	})

	t.Run("Test FullBucket", func(t *testing.T) {

		for i := 0; i < bucketSize-3; i++ {
			bucket.AddContact(NewContact(NewRandomKademliaID(), ""))
		}
		assert.Equal(t, bucket.Len(), bucketSize)

		assert.Equal(t, bucket.list.Back().Value.(Contact), contact1)

		newContact := NewContact(NewRandomKademliaID(), "")
		bucket.AddContact(newContact)
		assert.Equal(t, bucket.Len(), bucketSize)
		assert.Equal(t, bucket.list.Front().Value.(Contact), newContact)
	})
}
