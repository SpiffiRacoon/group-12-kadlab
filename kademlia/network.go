package kademlia

type Network struct {
	RoutingTable RoutingTable
	Me 	   Contact
}

type Message struct {
	MsgType string
	Body   string
	Sender Contact
}

func NewNetwork(me Contact) *Network {
	network := &Network{}
	network.Me = me
	network.RoutingTable = *NewRoutingTable(me)
	return network
}

func (network *Network) Listen(ip string, port int) error{
	// TODO
	return nil
}

func (network *Network) SendPingMessage(contact *Contact) Message{
	// TODO
	return Message{"test", "test", *contact}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
