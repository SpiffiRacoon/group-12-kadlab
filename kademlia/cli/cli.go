package cli

import (
	"d7024e/kademlia"
	"fmt"
	"os"
	"strings"
)

func Kcli(input string, node *kademlia.Kademlia) {
	trimmedinput := strings.TrimSpace(input)
	commandNdata := strings.Fields(trimmedinput)

	if len(commandNdata) == 0 {
		fmt.Println("No command given.")
	} else {
		command := commandNdata[0]
		switch command {
		case "ping":
			if len(commandNdata) == 2 {
				id := kademlia.NewKademliaID(commandNdata[1])
				contact, err := node.LookupContact(id)
				if err != nil {
					fmt.Println("failed to fetch contact from target id:", err)
				} else {
					node.Network.SendPingMessage(&contact[0])
				}
			} else {
				fmt.Println("This command needs one additional argument")
			}
		case "put":
			if len(commandNdata) == 2 {
				node.Store([]byte(commandNdata[1]))
			} else {
				fmt.Println("This command needs one additional argument")
			}
		case "get":
			if len(commandNdata) == 2 {
				recvData := node.LookupData(commandNdata[1])
				print(recvData)
			} else {
				fmt.Println("This comman needs one additional arguments")
			}
		case "print":
			node.Network.RoutingTable.PrintRoutingTable()
		case "exit":
			fmt.Println("shutting down node...")
			os.Exit(0)
		default:
			fmt.Println("not a valid argument")

		}
		fmt.Print("kCLI@", node.Me.ID, " % ")
	}
}
