package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"sort"
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
	kademlia.Network.RoutingTable.PrintRoutingTable()

	time.Sleep(20 * time.Second)
	kademlia.Network.RoutingTable.PrintRoutingTable()
}

func (kademlia *Kademlia) PopulateNetwork() {
	fmt.Println("Populating the network...")

	queriedNodes := []Contact{kademlia.Me}

	for i := 0; i < IDLength*8; i = 2*i + 1 { //Results in 8 iterations
		id := kademlia.Network.RoutingTable.GenerateIDForBucket(i)
		contacts, err := kademlia.LookupContact(id)
		if err != nil {
			fmt.Println("Error finding contacts during network population")
			continue
		}

		for _, contact := range contacts {
			if containsContact(queriedNodes, contact) {
				continue
			}

			queriedNodes = append(queriedNodes, contact)

			err = kademlia.Network.SendPingMessage(&contact)
			if err != nil {
				fmt.Println("Error pinging contact, node may be down")
				continue
			}

			kademlia.Network.RoutingTable.AddContact(contact)
			err = kademlia.Network.SendJoinMessage(&contact)
			if err != nil {
				fmt.Println("Error sending join message")
			}
		}
	}

	fmt.Println("Network population complete.")
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	k := 3     // Number of closest nodes to query (bucket size)
	alpha := 3 // Degree of parallelism (number of nodes to query in each iteration)
	closestNodes := kademlia.Network.RoutingTable.FindClosestContacts(target, k)

	queriedNodes := []Contact{}
	foundCloser := true

	for foundCloser {
		foundCloser = false

		for i := 0; i < len(closestNodes) && i < alpha; i++ {
			node := closestNodes[i]

			if containsContact(queriedNodes, node) {
				continue
			}

			queriedNodes = append(queriedNodes, node)

			newNodes, err := kademlia.Network.SendFindContactMessage(&node, target)
			if err != nil {
				closestNodes = removeFromList(closestNodes, node)
				continue
			}

			closestNodes = append(closestNodes, newNodes...)

			sort.Slice(closestNodes, func(i, j int) bool {
				return closestNodes[i].ID.CalcDistance(target).Less(closestNodes[j].ID.CalcDistance(target))
			})

			if len(closestNodes) > k {
				closestNodes = closestNodes[:k]
			}

			foundCloser = true
		}
	}

	return closestNodes, nil
}

func containsContact(list []Contact, current Contact) bool {
	for _, i := range list {
		if i.ID.Equals(current.ID) {
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

func (kademlia *Kademlia) LookupData(key string) ([]byte, bool) {
	k := 3     // Number of closest nodes to query (bucket size)
	alpha := 3 // Degree of parallelism (number of nodes to query in each iteration)
	if kademlia.DataStorage[key] != nil {
		return kademlia.DataStorage[key], true
	}
	closestContacts := kademlia.Network.RoutingTable.FindClosestContacts(kademlia.Network.RoutingTable.me.ID, k)
	var queriedContacts []Contact
	foundCloser := true

	for foundCloser {
		foundCloser = false
		for i := 0; i < len(closestContacts) && i < alpha; i++ {
			contact := closestContacts[i]
			//if new node, add it to list of queried nodes
			if containsContact(queriedContacts, contact) {
				continue
			}
			queriedContacts = append(queriedContacts, contact)
			value, newContacts := kademlia.Network.SendFindDataMessage(key, &contact)
			if value != "" {
				return []byte(value), true
			}
			closestContacts = append(closestContacts, newContacts...)
			sort.Slice(closestContacts, func(i, j int) bool {
				return closestContacts[i].ID.CalcDistance(kademlia.Network.RoutingTable.me.ID).Less(closestContacts[j].ID.CalcDistance(kademlia.Network.RoutingTable.me.ID))
			})
			if len(closestContacts) > k {
				closestContacts = closestContacts[:k]
			}
			foundCloser = true
		}
	}
	//Case: key not found
	return nil, false
}

func (kademlia *Kademlia) ExtractData(hash string) (data []byte, exists bool) {
	val, exists := kademlia.DataStorage[hash]
	return val, exists
}

func (kademlia *Kademlia) Store(data []byte) (string, error) {
	key := kademlia.MakeKey(data)
	location := NewKademliaID(key)
	contacts, _ := kademlia.LookupContact(location)
	if len(contacts) == 0 {
		return "", fmt.Errorf("no contacts found")
	}
	for _, contact := range contacts {
		fmt.Println("Storing data at: ", location.String(), " on node: ", contact.Address)
		err := kademlia.Network.SendStoreMessage(data, key, &contact)
		if err != nil {
			return "", err
		}
	}
	return key, nil
}

func (kademlia *Kademlia) LocalStorage(data []byte, key string) {
	kademlia.DataStorage[key] = data
}

func (kademlia *Kademlia) MakeKey(value []byte) string {
	sha1 := sha1.Sum(value) //hashes the data
	key := hex.EncodeToString(sha1[:])
	return key
}
