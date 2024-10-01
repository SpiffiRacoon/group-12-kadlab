package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	// Create a new bucket
	bucket := newBucket()

	// Create a new contact
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:8000")
	contact3 := NewContact(NewKademliaID("22222222000000000000000000000000000000000"), "localhost:8000")

	t.Run("Test AddContact", func(t *testing.T) {
		// Add a contact to the bucket
		bucket.AddContact(contact1)

		// Check if the contact has been added to the bucket
		if bucket.Len() != 1 {
			t.Fatalf("Expected 1 contact but instead got %d", bucket.Len())
		}

		// Add another contact to the bucket
		bucket.AddContact(contact2)
		bucket.AddContact(contact3)
		if bucket.Len() != 3 {
			t.Fatalf("Expected 3 contact but instead got %d", bucket.Len())
		}
	})

	t.Run("Test GetContactAndCalcDistance", func(t *testing.T) {
		// Get a contact from the bucket
		contacts := bucket.GetContactAndCalcDistance(NewKademliaID("0000000000000000000000000000000000000000"))

		assert.Equal(t, 3, len(contacts))

		assert.Equal(t, contacts[2].distance, contact1.ID)
		assert.Equal(t, contacts[1].distance, contact2.ID)
		assert.Equal(t, contacts[0].distance, contact3.ID)
	})
}
