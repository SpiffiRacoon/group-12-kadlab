package kademlia

import (
	"fmt"
	"time"
)

type Kademlia struct {
	Me Contact
	BootstrapNode Contact
	Network Network
	isBootstrap bool
}

func NewKademlia(me Contact, isBootstrap bool) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Me = me
	kademlia.Network= *NewNetwork(me)
	kademlia.isBootstrap = isBootstrap
	return kademlia
}

func (kademlia *Kademlia) Start() {
	if !kademlia.isBootstrap {
		go func() {
			kademlia.JoinNetwork(&kademlia.BootstrapNode)
		}()
	}
	err := kademlia.Network.Listen("0.0.0.0", 3000)
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	fmt.Println("Joining network...")
	time.Sleep(2 * time.Second)

	kademlia.Network.RoutingTable.AddContact(*knownNode)
	kademlia.Network.SendJoinMessage(knownNode)
	contacts, err := kademlia.LookupContact(&kademlia.Me)
	if err != nil {
		fmt.Println("Error finding contacts")
		return
	}

	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
	}
	
	fmt.Println("Network joined.")
	kademlia.PopulateNetwork()
	fmt.Printf("kademlia.Network.RoutingTable.buckets: %v\n", kademlia.Network.RoutingTable.buckets)
}

func (kademlia *Kademlia) PopulateNetwork() {
    fmt.Println("Populating the network...")

    // Define the number of random IDs to search for and populate the network with
    numLookups := 10 

    for i := 0; i < numLookups; i++ {
        randomID := NewRandomKademliaID()

        contacts, err := kademlia.LookupContact(&Contact{ID: randomID})
        if err != nil {
            fmt.Println("Error finding contacts during network population")
            continue
        }

        for _, contact := range contacts {
            kademlia.Network.RoutingTable.AddContact(contact)
        }
    }

    fmt.Println("Network population complete.")
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	k := 3
	closestNodes := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, k)

	queriedNodes := []Contact{}


	for _, node := range closestNodes {
		if containsContact(queriedNodes, node) {
			break
		}

		queriedNodes = append(queriedNodes, node)

		newNodes, err := kademlia.Network.SendFindContactMessage(&node, target.ID)
		if err != nil {
			// If the node is not responding, remove it from our closestNodes list
			closestNodes = removeFromList(closestNodes, node)
			continue
		}

		closestNodes = append(closestNodes, newNodes...)

		if len(closestNodes) > k {
			closestNodes = closestNodes[:k]
		}
	}

	return closestNodes, nil
}

func containsContact(list []Contact, current Contact) bool {
	for _, i := range list {
		if i.ID == current.ID {
			return true
		}
	}
	return false
}

func removeFromList(list []Contact, current Contact) []Contact {
	for i, contact := range list {
		if contact.ID == current.ID {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
