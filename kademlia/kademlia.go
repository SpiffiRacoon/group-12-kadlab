package kademlia

import (
	"fmt"
	"math/big"
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
		//kademlia.Network.SendJoinMessage(&contact)
	}
	
	fmt.Println("Network joined.")
	kademlia.PopulateNetwork()
	fmt.Println("Network populated.")
	fmt.Printf("kademlia.Network.RoutingTable.buckets: %v\n", kademlia.Network.RoutingTable.buckets)
}

func (kademlia *Kademlia) PopulateNetwork() {
    fmt.Println("Populating the network...")

    // Define the number of random IDs to search for and populate the network with
    numLookups := 5 // You can adjust this based on how much you want to populate the network

    for i := 0; i < numLookups; i++ {
        // Generate a random node ID for network discovery
        randomID := NewRandomKademliaID()

        // Perform a lookup for the random node ID
        contacts, err := kademlia.LookupContact(&Contact{ID: randomID})
        if err != nil {
            fmt.Println("Error finding contacts during network population")
            continue
        }

        // Add the discovered contacts to the routing table
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


	// Recursive querying loop
	for _, node := range closestNodes {
		if containsContact(queriedNodes, node) {
			break
		}

		// Mark the node as queried
		queriedNodes = append(queriedNodes, node)

		// Perform a network lookup on this node (sending a FIND_NODE request)
		newNodes, err := kademlia.Network.SendFindContactMessage(&node, target.ID)
		if err != nil {
			// If the node is not responding, remove it from our closestNodes list
			closestNodes = removeFromList(closestNodes, node)
			continue
		}

		// Add new nodes to our closestNodes list
		closestNodes = append(closestNodes, newNodes...)

		// Sort nodes by their distance to targetID again after adding new nodes
		//closestNodes.Sort()

		// Keep only the closest k nodes
		if len(closestNodes) > k {
			closestNodes = closestNodes[:k]
		}
	}

	fmt.Println("Found ", len(closestNodes), " closest nodes")
	fmt.Println("Nodes: ", closestNodes)
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

func xorDistance(id1, id2 []byte) *big.Int {
	id1Big := new(big.Int).SetBytes(id1)
	id2Big := new(big.Int).SetBytes(id2)
	return new(big.Int).Xor(id1Big, id2Big)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
