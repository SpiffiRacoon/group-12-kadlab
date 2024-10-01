package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKademliaId(t *testing.T) {
	// Create a new KademliaID

	t.Run("Test NewKademliaID", func(t *testing.T) {
		// Testing equality of two KademliaIDs
		id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		id2 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		assert.Equal(t, id1, id2)

		// Testing inequality of two KademliaIDs
		id3 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		id4 := NewKademliaID("1111111100000000000000000000000000000000")
		assert.NotEqual(t, id3, id4)
	})

	t.Run("Test NewRandomKademliaID", func(t *testing.T) {
		// Testing inequality of two random KademliaIDs
		id1 := NewRandomKademliaID()
		id2 := NewRandomKademliaID()
		assert.NotEqual(t, id1, id2)
	})

	t.Run("Test Less", func(t *testing.T) {
		// Testing Less method of KademliaID
		id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		id2 := NewKademliaID("1111111100000000000000000000000000000000")
		assert.True(t, id2.Less(id1))
		assert.False(t, id1.Less(id2))
	})

	t.Run("Test Equals", func(t *testing.T) {
		// Testing Equals method of KademliaID
		id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		id2 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		assert.True(t, id1.Equals(id2))

		id3 := NewKademliaID("1111111100000000000000000000000000000000")
		assert.False(t, id1.Equals(id3))
	})

	t.Run("Test CalcDistance", func(t *testing.T) {
		// Testing CalcDistance method of KademliaID
		id1 := NewKademliaID("F000000000000000000000000000000000000000")
		id2 := NewKademliaID("1000000000000000000000000000000000000000")
		assert.Equal(t, id1.CalcDistance(id2), NewKademliaID("E000000000000000000000000000000000000000"))
		
		id3 := NewKademliaID("0000000000000000000000000000000000000000")
		assert.Equal(t, id1.CalcDistance(id1), id3)
	})

	t.Run("Test String", func(t *testing.T) {
		// Testing String method of KademliaID
		id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
		// OBS: String() returns lower case
		assert.Equal(t, id1.String(), "ffffffff00000000000000000000000000000000")

	})
}