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
	contacts := kademlia.Network.RoutingTable.FindClosestContacts(target.ID, 3)
	return contacts, nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
