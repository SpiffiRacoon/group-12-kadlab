package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

type Kademlia struct {
	Me            Contact
	BootstrapNode Contact
	Network       Network
	isBootstrap   bool
	dataStorage   map[string][]byte
}

func NewKademlia(me Contact, isBootstrap bool) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Me = me
	kademlia.Network = *NewNetwork(me)
	kademlia.isBootstrap = isBootstrap
	return kademlia
}

func (kademlia *Kademlia) Start() {
	if !kademlia.isBootstrap {
		go func() {
			kademlia.JoinNetwork()
		}()
	}
	err := kademlia.Network.Listen("0.0.0.0", 3000)
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("Joining network...")
	time.Sleep(2 * time.Second)

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

// changed target from *Contact to *KademliaID so it can go straight as input to "FindclosestContacts"
func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	//	probedContacts := new([]Contact)                                                           //a list of already visited contacts
	var closestContacts *[]Contact // this var holds a pointer to a list
	//	currClosest := NewContact((NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")), "") //Sets the current closest contact as a 160-bit KademliaID of all ones, sets data as empty string
	//	currClosest.distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")           //distance becomes
	//
	//	alphaClosestContacts := kademlia.Network.RoutingTable.FindClosestContacts(target, alpha) //alpha is a system-wide concurrency parameter, such as 3(page 6 of research paper)
	//	closestContacts = &alphaClosestContacts
	//	//TODO everything

	return *closestContacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	sha1 := sha1.Sum(data) //hashes the data
	key := hex.EncodeToString(sha1[:])
	//contact := kademlia.LookupContact(NewKademliaID((key)))   //TODO implement LookupContact
	//
	go kademlia.Network.SendStoreMessage(data, key, &kademlia.BootstrapNode) //TODO implement SendStoreMessage, should send store message to bootstrap node

}

func (kademlia *Kademlia) LocalStorage(data []byte, key string) {
	//TODO lösa hur key value paret lagras lokalt på en nods dataStorage
}
