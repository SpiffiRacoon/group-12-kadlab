package kademlia

import (
	"container/list"
)

// bucket definition
// contains a List
type bucket struct {
	list    *list.List
	network *Network
}

// newBucket returns a new instance of a bucket
func newBucket(network *Network) *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.network = network
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			oldest := bucket.list.Back().Value.(Contact)
			err := bucket.network.SendPingMessage(&oldest)
			if err != nil {
				bucket.list.Remove(bucket.list.Back())
				bucket.list.PushFront(contact)
			}
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
