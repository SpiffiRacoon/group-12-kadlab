package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContacts(t *testing.T) {
	// Create a new contact
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000")
	contact3 := NewContact(NewKademliaID("2222222200000000000000000000000000000000"), "localhost:8000")

	t.Run("Test CalcDistance", func(t *testing.T) {
		contact1.CalcDistance(contact2.ID)
		assert.Equal(t, contact1.distance, NewKademliaID("EEEEEEEE00000000000000000000000000000000"))

		contact2.CalcDistance(contact3.ID)
		assert.Equal(t, contact2.distance, NewKademliaID("3333333300000000000000000000000000000000"))
	})

	t.Run("Test Less", func(t *testing.T) {
		assert.False(t, contact1.Less(&contact2))
		assert.True(t, contact2.Less(&contact1))
	})

	t.Run("Test String", func(t *testing.T) {
		assert.Equal(t, contact1.String(), `contact("ffffffff00000000000000000000000000000000", "localhost:8000")`)
		assert.Equal(t, contact2.String(), `contact("1111111100000000000000000000000000000000", "localhost:8000")`)
		assert.Equal(t, contact3.String(), `contact("2222222200000000000000000000000000000000", "localhost:8000")`)
	})
}