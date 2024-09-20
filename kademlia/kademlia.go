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

	/*
	ping := kademlia.Network.SendPingMessage(&kademlia.BootstrapNode)
	if !ping {
		fmt.Println("Bootstrap node not responding")
		return
	}

	contacts, err := kademlia.Network.SendFindContactMessage(&kademlia.BootstrapNode)
	if err != nil {
		fmt.Println("Error finding contacts")
		return
	}
*/
	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
	}	

	// TODO
	// check if bootstrap node is alive
	// send ping message to bootstrap node
	// if response, add bootstrap node to routing table
	// send find contact message to bootstrap node
	// if response, add contacts to routing table

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


	/*
	targetID := target.ID
	alpha := 3
	checkedContacts := new([]Contact)
	var closestList *[]Contact
	alphaclosestList := kademlia.Network.RoutingTable.FindClosestContacts(targetID, alpha)
	closestList = &alphaclosestList

	currentClosest := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "")
	currentClosest.distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	for {
		updateClosest := false
		numChecked := 0
		for i := 0; i < len(*closestList) && numChecked < alpha; i++ {
			if containsContact(*checkedContacts, (*closestList)[i]) {
				continue
			} else {
				templist, err := kademlia.Network.SendFindContactMessage(&(*closestList)[i], *targetID)
				if err != nil {
					//kademlia.Network.RoutingTable.RemoveContact((*closestList)[i])
					*closestList = removeFromList(*closestList, (*closestList)[i])
					continue
				} else {
					*checkedContacts = append(*checkedContacts, (*closestList)[i])
					bucket := kademlia.Network.RoutingTable.buckets[kademlia.Network.RoutingTable.getBucketIndex((*closestList)[i].ID)]
					// if there is space in the bucket add the node
					kademlia.updateBucket(*bucket, (*closestList)[i])
					// append contacts to shortlist if err is none
					for i := 0; i < len(templist); i++ {
						templist[i].CalcDistance(targetID)
					}
					*closestList, currentClosest, updateClosest = kademlia.addUniqueContacts(templist, *closestList, currentClosest, updateClosest)
					numChecked++
				}

			}
		}
		if !updateClosest || len(*checkedContacts) >= 20 {
			break
		}
	}

	return *closestList, nil
	*/
	/*
	
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, 3)
	fmt.Println("Num found contacts: ", len(contacts))
	return contacts, nil
	*/
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
