package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Kademlia struct {
	Me            Contact
	BootstrapNode Contact
	Network       Network
	IsBootstrap   bool
	DataStorage   map[string][]byte
}

func NewKademlia(me Contact, isBootstrap bool) *Kademlia {
	kademlia := Kademlia{}
	kademlia.Me = me
	kademlia.Network = *NewNetwork(me, &kademlia)
	kademlia.IsBootstrap = isBootstrap
	kademlia.DataStorage = make(map[string][]byte)
	return &kademlia
}

func (kademlia *Kademlia) Start() {
	if !kademlia.IsBootstrap {
		go func() {
			kademlia.JoinNetwork(&kademlia.BootstrapNode)
		}()
	}

	host, portStr, err := net.SplitHostPort(kademlia.Me.Address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Error converting port to int:", err)
		return
	}
	err = kademlia.Network.Listen(host, port)
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	fmt.Println("Joining network...")
	time.Sleep(2 * time.Second)

	kademlia.Network.RoutingTable.AddContact(*knownNode)
	kademlia.Network.SendJoinMessage(knownNode)
	contacts, err := kademlia.LookupContact(kademlia.Me.ID)
	if err != nil {
		fmt.Println("Error finding contacts")
		return
	}

	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
	}

	fmt.Println("Network joined.")
	kademlia.PopulateNetwork()
}

func (kademlia *Kademlia) PopulateNetwork() {
	fmt.Println("Populating the network...")

	// Define the number of random IDs to search for and populate the network with
	numLookups := 10

	for i := 0; i < numLookups; i++ {
		randomID := NewRandomKademliaID()

		contacts, err := kademlia.LookupContact(randomID)
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

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	k := 3
	closestNodes := kademlia.Network.RoutingTable.FindClosestContacts(target, k)

	queriedNodes := []Contact{}

	for _, node := range closestNodes {
		if containsContact(queriedNodes, node) {
			break
		}

		queriedNodes = append(queriedNodes, node)

		newNodes, err := kademlia.Network.SendFindContactMessage(&node, target)
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

	//TODO: sort the contacts?

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

func (kademlia *Kademlia) LookupData(key string) []byte {
	if kademlia.DataStorage[key] != nil {
		return kademlia.DataStorage[key]
	}
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(kademlia.Network.RoutingTable.me.ID, 5)
	var queriedContacts []Contact

	for _, contact := range contacts {
		//if new node, add it to list of queried nodes
		if !containsContact(queriedContacts, contact) {
			queriedContacts = append(queriedContacts, contact)
			response, newNodes := kademlia.Network.SendFindDataMessage(key, &contact)
			//If response is empty then it has gotten suggested nodes from SendFindDataMessage which should be queried in later iterations
			if response == "" {
				contacts = append(contacts, newNodes...)
			} else {
				//if the response is not empty, the key has been found
				return []byte(response)
			}
		}

	}
	//Case: key not found
	return nil
}

func (kademlia *Kademlia) ExtractData(hash string) (data []byte, exists bool) {
	val, exists := kademlia.DataStorage[hash]
	return val, exists
}

func (kademlia *Kademlia) Store(data []byte) error {
	key := kademlia.MakeKey(data)
	location := NewKademliaID(key)
	contacts, _ := kademlia.LookupContact(location)

	if len(contacts) == 0 {
		return fmt.Errorf("no contacts found for key %s", key)
	}
	for _, contact := range contacts {
		kademlia.Network.SendStoreMessage(data, key, &contact)
	}
	return nil
}

func (kademlia *Kademlia) LocalStorage(data []byte, key string) {
	kademlia.DataStorage[key] = data
}

func (kademlia *Kademlia) MakeKey(value []byte) string {
	sha1 := sha1.Sum(value) //hashes the data
	key := hex.EncodeToString(sha1[:])
	return key
}
