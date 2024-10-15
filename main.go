package main

import (
	"bufio"
	"d7024e/kademlia"
	"d7024e/kademlia/cli"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func docker_health_check(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "I'm healthy!")
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	local_NET_IP, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)
	}

	localIP := local_NET_IP[0].String()

	NodePort := os.Getenv("NODE_PORT")

	BOOTSTRAP_HOSTNAME := os.Getenv("BOOSTRAP_NODE_HOSTNAME")
	IS_BOOTSTRAP := os.Getenv("IS_BOOTSTRAP") == "true"

	var bootstrapIP string

	if !IS_BOOTSTRAP {
		bootstrap_NET_IP, err := net.LookupIP(BOOTSTRAP_HOSTNAME)
		if err != nil {
			panic(err)
		}
		bootstrapIP = bootstrap_NET_IP[0].String()
	}
	
	contact := kademlia.NewContact(kademlia.NewRandomKademliaID(), localIP+":"+NodePort)
	fmt.Println("Me: " + contact.String())

	node := kademlia.NewKademlia(contact, IS_BOOTSTRAP)

	if IS_BOOTSTRAP {
		http.HandleFunc("/health", docker_health_check)
		go http.ListenAndServe("127.0.0.1:80", nil)
	}

	bootsrapNode := kademlia.NewContact(
		kademlia.NewKademliaID("B0075712A9000000000000000000000000000000"), 
		bootstrapIP+":3000")

	node.BootstrapNode = bootsrapNode
	go node.Start()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		input := scanner.Text()
		cli.Kcli(input, node)
	}

}
