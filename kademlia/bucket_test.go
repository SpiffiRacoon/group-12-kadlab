package kademlia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	// Create a new bucket
	//contact0 := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000")
	bucket := newBucket(nil)

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

		for i := 0; i < bucketSize; i++ {
			bucket.AddContact(NewContact(NewRandomKademliaID(), ""))
		}
		fmt.Println(bucket.Len())
		
	})
}
