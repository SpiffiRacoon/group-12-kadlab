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
	kademlia := &Kademlia{}
	kademlia.Me = me
	kademlia.Network = *NewNetwork(me)
	kademlia.IsBootstrap = isBootstrap
	return kademlia
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

	// Define the number of random IDs to search for and populate the network with
	for i := 0; i < IDLength*8; i++ { //2*i+1 Results in 8 iterations
		id := kademlia.Network.RoutingTable.GenerateIDForBucket(i)
		contacts, err := kademlia.LookupContact(id)
		if err != nil {
			fmt.Println("Error finding contacts during network population")
			continue
		}

		for _, contact := range contacts { //TODO: still contacts same mode multiple times?
			if containsContact(queriedNodes, contact) {
				continue
			}
			queriedNodes = append(queriedNodes, contact)
			err = kademlia.Network.SendPingMessage(&contact)
			if err != nil {
				fmt.Println("Error pinging contact during network population, node may be down")
			}else {
				kademlia.Network.RoutingTable.AddContact(contact)
				kademlia.Network.SendJoinMessage(&contact)
			}
		}
		fmt.Printf("Found %d contacts for Bucket %d\n", len(contacts), i)
	}

	fmt.Println("Network population complete.")
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	k := 3  // Number of closest nodes to query (bucket size)
	alpha := 3  // Degree of parallelism (number of nodes to query in each iteration)
	closestNodes := kademlia.Network.RoutingTable.FindClosestContacts(target, k)

	queriedNodes := []Contact{}  // Tracks nodes that have already been queried
	foundCloser := true  // Tracks if new, closer nodes are found

	for foundCloser {
		foundCloser = false  // Reset the flag; will set to true if closer nodes are found

		// Iterate over the `closestNodes` and query unqueried nodes
		for i := 0; i < len(closestNodes) && i < alpha; i++ {
			node := closestNodes[i]

			// Skip nodes that have already been queried
			if containsContact(queriedNodes, node) {
				continue
			}

			queriedNodes = append(queriedNodes, node)

			// Query the node for closer contacts
			newNodes, err := kademlia.Network.SendFindContactMessage(&node, target)
			if err != nil {
				// Remove non-responsive nodes from the closestNodes list
				closestNodes = removeFromList(closestNodes, node)
				continue
			}

			// Add the newly discovered nodes to the list of closest nodes
			closestNodes = append(closestNodes, newNodes...)

			// Sort by proximity to the target (closer nodes first)
			sort.Slice(closestNodes, func(i, j int) bool {
				return closestNodes[i].ID.CalcDistance(target).Less(closestNodes[j].ID.CalcDistance(target))
			})

			// Ensure we only keep the k closest nodes
			if len(closestNodes) > k {
				closestNodes = closestNodes[:k]
			}

			// If new closer nodes were added, we continue querying
			foundCloser = true
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

func (kademlia *Kademlia) LookupData(hash string) string {
	data := kademlia.ExtractData(hash)
	if data != nil {
		return string(data)
	}
	location := NewKademliaID(hash)
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(location, 5)
	for _, contact := range contacts {
		searches, found := kademlia.Network.SendFindDataMessage(hash, &contact)
		if found == nil {
			return string(kademlia.ExtractData(searches)) //Unclear if this is the correct way, want to extract the data stored in node with kademliaID "contact"
		}
	}
	return "Did not find the data in the closest contacts" //not sure if this is what we want either
}

func (kademlia *Kademlia) ExtractData(hash string) (data []byte) {
	res := kademlia.DataStorage[hash]
	return res
}

func (kademlia *Kademlia) Store(data []byte) {
	sha1 := sha1.Sum(data) //hashes the data
	key := hex.EncodeToString(sha1[:])
	location := NewKademliaID(key)
	contacts, _ := kademlia.LookupContact(location)

	if len(contacts) <= 0 {
		fmt.Println("Error, no suitable nodes to store the data could be found")
	} else {
		//blank because we do not care about the iteration variable
		//we basically do a for each contact in contacts
		for _, contact := range contacts {
			kademlia.Network.SendStoreMessage(data, key, &contact)
		}
		fmt.Println(string(data) + " is now stored in the key: " + key)
	}
}

func (kademlia *Kademlia) LocalStorage(data []byte, key string) {
	kademlia.DataStorage[key] = data
}
